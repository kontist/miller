#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "mlr_dsl_cst.h"

static sllv_t* mlr_dsl_cst_alloc_from_statement_list(sllv_t* pasts, int type_inferencing);
static mlr_dsl_cst_statement_t* cst_statement_alloc(mlr_dsl_ast_node_t* past, int type_inferencing);
static void cst_statement_free(mlr_dsl_cst_statement_t* pstatement);

static mlr_dsl_cst_statement_item_t* mlr_dsl_cst_statement_item_alloc(
	int lhs_type,
	char* output_field_name,
	sllv_t* pmoosvar_lhs_keylist_evaluators,
	lrec_evaluator_t* prhs_evaluator);
static void cst_statement_item_free(mlr_dsl_cst_statement_item_t* pitem);

// ----------------------------------------------------------------
// At present (initial oosvar experimens, January 2016) the begin/main/end are organized as follows:
//
// Input:
//
//   mlr put 'begin @a = 1; begin @b = 0; @b = $x + @a; end emit @a, @b'
//
// i.e. there are separate begin keywords one per statement, rather than
//
//   mlr put 'begin { @a = 1; @b = 0}; @b = $x + @a; end { emit @a, @b }'
//
// Corresponding list of per-statement ASTs:
//   begin (begin):
//       = (oosvar_assignment):
//           a (oosvar_name).
//           1 (strnum_literal).
//   begin (begin):
//       = (oosvar_assignment):
//           b (oosvar_name).
//           0 (strnum_literal).
//   = (oosvar_assignment):
//       b (oosvar_name).
//       + (operator):
//           x (field_name).
//           a (oosvar_name).
//   end (end):
//       emit (emit):
//           a (oosvar_name).
//           b (oosvar_name).
//
// (Note that the AST input is a list of per-statement ASTs, rather than a single root-node AST with per-statement child
// nodes.)
//
// So our job here at present is to loop through the per-statement ASTs, splitting them out by begin/main/end.

mlr_dsl_cst_t* mlr_dsl_cst_alloc(mlr_dsl_ast_t* past, int type_inferencing) {
	mlr_dsl_cst_t* pcst = mlr_malloc_or_die(sizeof(mlr_dsl_cst_t));

	pcst->pbegin_statements = mlr_dsl_cst_alloc_from_statement_list(past->pbegin_statements, type_inferencing);
	pcst->pmain_statements  = mlr_dsl_cst_alloc_from_statement_list(past->pmain_statements,  type_inferencing);
	pcst->pend_statements   = mlr_dsl_cst_alloc_from_statement_list(past->pend_statements,   type_inferencing);

	return pcst;
}

static sllv_t* mlr_dsl_cst_alloc_from_statement_list(sllv_t* pasts, int type_inferencing) {
	sllv_t* pstatements = sllv_alloc();
	for (sllve_t* pe = pasts->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* past = pe->pvvalue;
		mlr_dsl_cst_statement_t* pstatement = cst_statement_alloc(past, type_inferencing);
		sllv_add(pstatements, pstatement);
	}
	return pstatements;
}

void mlr_dsl_cst_free(mlr_dsl_cst_t* pcst) {
	if (pcst == NULL)
		return;
	for (sllve_t* pe = pcst->pbegin_statements->phead; pe != NULL; pe = pe->pnext)
		cst_statement_free(pe->pvvalue);
	for (sllve_t* pe = pcst->pmain_statements->phead; pe != NULL; pe = pe->pnext)
		cst_statement_free(pe->pvvalue);
	for (sllve_t* pe = pcst->pend_statements->phead; pe != NULL; pe = pe->pnext)
		cst_statement_free(pe->pvvalue);
	sllv_free(pcst->pbegin_statements);
	sllv_free(pcst->pmain_statements);
	sllv_free(pcst->pend_statements);
	free(pcst);
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* cst_statement_alloc(mlr_dsl_ast_node_t* past, int type_inferencing) {
	mlr_dsl_cst_statement_t* pstatement = mlr_malloc_or_die(sizeof(mlr_dsl_cst_statement_t));

	pstatement->ast_node_type = past->type;
	pstatement->pitems = sllv_alloc();

	if (past->type == MD_AST_NODE_TYPE_SREC_ASSIGNMENT) {
		if ((past->pchildren == NULL) || (past->pchildren->length != 2)) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}

		mlr_dsl_ast_node_t* pleft  = past->pchildren->phead->pvvalue;
		mlr_dsl_ast_node_t* pright = past->pchildren->phead->pnext->pvvalue;

		if (pleft->type != MD_AST_NODE_TYPE_FIELD_NAME) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		} else if (pleft->pchildren != NULL) {
			fprintf(stderr, "%s: coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}

		sllv_add(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			MLR_DSL_CST_LHS_TYPE_SREC,
			pleft->text,
			NULL,
			lrec_evaluator_alloc_from_ast(pright, type_inferencing)));

	} else if (past->type == MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT) {
		if ((past->pchildren == NULL) || (past->pchildren->length != 2)) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}

		mlr_dsl_ast_node_t* pleft  = past->pchildren->phead->pvvalue;
		mlr_dsl_ast_node_t* pright = past->pchildren->phead->pnext->pvvalue;

		if (pleft->type != MD_AST_NODE_TYPE_OOSVAR_NAME) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		} else if (pleft->pchildren != NULL) {
			fprintf(stderr, "%s: coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}

		sllv_add(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			MLR_DSL_CST_LHS_TYPE_OOSVAR,
			pleft->text,
			NULL,
			lrec_evaluator_alloc_from_ast(pright, type_inferencing)));

	} else if (past->type == MD_AST_NODE_TYPE_MOOSVAR_ASSIGNMENT) {
		sllv_t* pmoosvar_lhs_keylist_evaluators = sllv_alloc();

		mlr_dsl_ast_node_t* pleft  = past->pchildren->phead->pvvalue;
		mlr_dsl_ast_node_t* pright = past->pchildren->phead->pnext->pvvalue;

		if (pleft->type != MD_AST_NODE_TYPE_MOOSVAR_NAME && pleft->type != MD_AST_NODE_TYPE_MOOSVAR_LEVEL_KEY) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}

		if (pleft->type == MD_AST_NODE_TYPE_MOOSVAR_NAME) {
			sllv_add(pmoosvar_lhs_keylist_evaluators,
				// xxx need a version with no regex-captures.
				lrec_evaluator_alloc_from_strnum_literal(mlr_strdup_or_die(pleft->text), TYPE_INFER_STRING_ONLY));
		} else {

			mlr_dsl_ast_node_t* pnode = pleft;
			while (TRUE) {
				// Example AST:
				// % mlr put -v 'begin{@@x[1]["2"][$3][@4]=5}' /dev/null
				// AST BEGIN STATEMENTS (1):
				// = (moosvar_assignment):
				//     [] (moosvar_level_key):
				//         [] (moosvar_level_key):
				//             [] (moosvar_level_key):
				//                 [] (moosvar_level_key):
				//                     x (moosvar_name).
				//                     1 (strnum_literal).
				//                 2 (strnum_literal).
				//             3 (field_name).
				//         4 (oosvar_name).
				//     5 (strnum_literal).
				//
				// Here past is the =; pright is the 5; pleft is the string of bracket references
				// ending at the moosvar name.

				// xxx rename pfoo
				if (pnode->type == MD_AST_NODE_TYPE_MOOSVAR_LEVEL_KEY) {
					mlr_dsl_ast_node_t* pfoo = pnode->pchildren->phead->pnext->pvvalue;
					sllv_add(pmoosvar_lhs_keylist_evaluators,
						lrec_evaluator_alloc_from_ast(pfoo, type_inferencing));
				} else {
					sllv_add(pmoosvar_lhs_keylist_evaluators,
						// xxx big comment here. this is confusing.
						lrec_evaluator_alloc_from_strnum_literal(mlr_strdup_or_die(pnode->text), TYPE_INFER_STRING_ONLY));
				}
				if (pnode->pchildren == NULL)
					break;
				pnode = pnode->pchildren->phead->pvvalue;
			}
			// Bracket operators come in from the right. So the highest AST node is the rightmost
			// map, and the lowest is the moosvar name.
			sllv_reverse(pmoosvar_lhs_keylist_evaluators);
			// xxx just make an sllv_add_at_head function
		}

		sllv_add(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			MLR_DSL_CST_LHS_TYPE_MOOSVAR,
			NULL,
			pmoosvar_lhs_keylist_evaluators,
			lrec_evaluator_alloc_from_ast(pright, type_inferencing)));

	} else if (past->type == MD_AST_NODE_TYPE_FILTER) {
		mlr_dsl_ast_node_t* pnode = past->pchildren->phead->pvvalue;
		sllv_add(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			MLR_DSL_CST_LHS_TYPE_OOSVAR,
			NULL,
			NULL,
			lrec_evaluator_alloc_from_ast(pnode, type_inferencing)));

	} else if (past->type == MD_AST_NODE_TYPE_GATE) {
		mlr_dsl_ast_node_t* pnode = past->pchildren->phead->pvvalue;
		sllv_add(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			MLR_DSL_CST_LHS_TYPE_OOSVAR,
			NULL,
			NULL,
			lrec_evaluator_alloc_from_ast(pnode, type_inferencing)));

	} else if (past->type == MD_AST_NODE_TYPE_EMIT) {
		// Loop over oosvar names to emit in e.g. 'emit @a, @b, @c'.
		for (sllve_t* pe = past->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pnode = pe->pvvalue;
			sllv_add(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
				MLR_DSL_CST_LHS_TYPE_OOSVAR,
				pnode->text,
				NULL,
				lrec_evaluator_alloc_from_ast(pnode, type_inferencing)));
		}

	} else if (past->type == MD_AST_NODE_TYPE_DUMP) {
		// No arguments: the node-type alone suffices for the caller to be able to execute this.

	} else { // Bare-boolean statement
		sllv_add(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			MLR_DSL_CST_LHS_TYPE_NONE,
			NULL,
			NULL,
			lrec_evaluator_alloc_from_ast(past, type_inferencing)));
	}

	return pstatement;
}

static void cst_statement_free(mlr_dsl_cst_statement_t* pstatement) {
	for (sllve_t* pe = pstatement->pitems->phead; pe != NULL; pe = pe->pnext)
		cst_statement_item_free(pe->pvvalue);
	sllv_free(pstatement->pitems);
	free(pstatement);
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_item_t* mlr_dsl_cst_statement_item_alloc(
	int lhs_type,
	char* output_field_name,
	sllv_t* pmoosvar_lhs_keylist_evaluators,
	lrec_evaluator_t* prhs_evaluator)
{
	mlr_dsl_cst_statement_item_t* pitem = mlr_malloc_or_die(sizeof(mlr_dsl_cst_statement_item_t));
	pitem->lhs_type = lhs_type;
	pitem->output_field_name = output_field_name == NULL ? NULL : mlr_strdup_or_die(output_field_name);
	pitem->pmoosvar_lhs_keylist_evaluators = pmoosvar_lhs_keylist_evaluators;
	pitem->prhs_evaluator = prhs_evaluator;
	return pitem;
}

static void cst_statement_item_free(mlr_dsl_cst_statement_item_t* pitem) {
	if (pitem == NULL)
		return;
	free(pitem->output_field_name);
	pitem->prhs_evaluator->pfree_func(pitem->prhs_evaluator);
	if (pitem->pmoosvar_lhs_keylist_evaluators != NULL) {
		for (sllve_t* pe = pitem->pmoosvar_lhs_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
			lrec_evaluator_t* pevaluator = pe->pvvalue;
			pevaluator->pfree_func(pevaluator);
		}
		sllv_free(pitem->pmoosvar_lhs_keylist_evaluators);
	}
	free(pitem);
}

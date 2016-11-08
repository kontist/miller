#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <ctype.h> // for tolower(), toupper()
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/mlrregex.h"
#include "lib/mtrand.h"
#include "mapping/mapper.h"
#include "mapping/rval_evaluators.h"
#include "mapping/function_manager.h"
#include "mapping/mlr_dsl_cst.h" // xxx only for allocate_keylist_evaluators_from_ast_node -- xxx move
#include "mapping/context_flags.h"

// ================================================================
// See comments in rval_evaluators.h
// ================================================================

// ----------------------------------------------------------------
rxval_evaluator_t* rxval_evaluator_alloc_from_ast(mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr,
	int type_inferencing, int context_flags)
{
	switch(pnode->type) {

	case MD_AST_NODE_TYPE_MAP_LITERAL:
		return NULL; // xxx XXX mapvar stub
		break;

	case MD_AST_NODE_TYPE_FUNCTION_CALLSITE:
		return NULL; // xxx XXX mapvar stub
		break;

	case MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE:
		return rxval_evaluator_alloc_from_nonindexed_local_variable(pnode, pfmgr, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_INDEXED_LOCAL_VARIABLE:
		return rxval_evaluator_alloc_from_indexed_local_variable(pnode, pfmgr, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_FULL_SREC:
		return NULL; // xxx XXX mapvar stub
		break;

	case MD_AST_NODE_TYPE_OOSVAR_KEYLIST:
		return NULL; // xxx XXX mapvar stub
		break;

	case MD_AST_NODE_TYPE_FULL_OOSVAR:
		return NULL; // xxx XXX mapvar stub
		break;

	default:
		return rxval_evaluator_alloc_wrapping_rval(pnode, pfmgr, type_inferencing, context_flags);
		break;
	}
}

// ================================================================
typedef struct _rxval_evaluator_from_nonindexed_local_variable_state_t {
	int vardef_frame_relative_index;
} rxval_evaluator_from_nonindexed_local_variable_state_t;

mlhmmv_value_t rxval_evaluator_from_nonindexed_local_variable_func(void* pvstate, variables_t* pvars) {
	rxval_evaluator_from_nonindexed_local_variable_state_t* pstate = pvstate;
	local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
	mlhmmv_value_t* pxval = local_stack_frame_get_map_value(pframe, pstate->vardef_frame_relative_index, NULL);
	return mlhmmv_copy_aux(pxval);
}

static void rxval_evaluator_from_nonindexed_local_variable_free(rxval_evaluator_t* prxval_evaluator) {
	rxval_evaluator_from_nonindexed_local_variable_state_t* pstate = prxval_evaluator->pvstate;
	free(pstate);
	free(prxval_evaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_from_nonindexed_local_variable(
	mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr, int type_inferencing, int context_flags)
{
	rxval_evaluator_from_nonindexed_local_variable_state_t* pstate = mlr_malloc_or_die(
		sizeof(rxval_evaluator_from_nonindexed_local_variable_state_t));
	MLR_INTERNAL_CODING_ERROR_IF(pnode->vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->vardef_frame_relative_index = pnode->vardef_frame_relative_index;

	rxval_evaluator_t* prxval_evaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	prxval_evaluator->pvstate       = pstate;
	prxval_evaluator->pprocess_func = rxval_evaluator_from_nonindexed_local_variable_func;
	prxval_evaluator->pfree_func    = rxval_evaluator_from_nonindexed_local_variable_free;

	return prxval_evaluator;
}

// ================================================================
typedef struct _rxval_evaluator_from_indexed_local_variable_state_t {
	int vardef_frame_relative_index;
	sllv_t* pkeylist_evaluators;
} rxval_evaluator_from_indexed_local_variable_state_t;

mlhmmv_value_t rxval_evaluator_from_indexed_local_variable_func(void* pvstate, variables_t* pvars) {
	rxval_evaluator_from_indexed_local_variable_state_t* pstate = pvstate;

	int all_non_null_or_error = TRUE;
	sllmv_t* pmvkeys = evaluate_list(pstate->pkeylist_evaluators, pvars, &all_non_null_or_error);

	if (all_non_null_or_error) {
		local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
		mlhmmv_value_t* pxval = local_stack_frame_get_map_value(pframe, pstate->vardef_frame_relative_index,
			pmvkeys);
		sllmv_free(pmvkeys);
		return mlhmmv_copy_aux(pxval);
	} else {
		sllmv_free(pmvkeys);
		return mlhmmv_value_transfer_terminal(mv_absent());
	}
}

static void rxval_evaluator_from_indexed_local_variable_free(rxval_evaluator_t* prxval_evaluator) {
	rxval_evaluator_from_indexed_local_variable_state_t* pstate = prxval_evaluator->pvstate;
	for (sllve_t* pe = pstate->pkeylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
		rval_evaluator_t* prval_evaluator = pe->pvvalue;
		prval_evaluator->pfree_func(prval_evaluator);
	}
	sllv_free(pstate->pkeylist_evaluators);
	free(pstate);
	free(prxval_evaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_from_indexed_local_variable(
	mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr, int type_inferencing, int context_flags)
{
	rxval_evaluator_from_indexed_local_variable_state_t* pstate = mlr_malloc_or_die(
		sizeof(rxval_evaluator_from_indexed_local_variable_state_t));
	MLR_INTERNAL_CODING_ERROR_IF(pnode->vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->vardef_frame_relative_index = pnode->vardef_frame_relative_index;
	pstate->pkeylist_evaluators = allocate_keylist_evaluators_from_ast_node(
		pnode, pfmgr, type_inferencing, context_flags);

	rxval_evaluator_t* prxval_evaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	prxval_evaluator->pvstate       = pstate;
	prxval_evaluator->pprocess_func = rxval_evaluator_from_indexed_local_variable_func;
	prxval_evaluator->pfree_func    = rxval_evaluator_from_indexed_local_variable_free;

	return prxval_evaluator;
}

// ================================================================
typedef struct _rxval_evaluator_wrapping_rval_state_t {
	rval_evaluator_t* prval_evaluator;
} rxval_evaluator_wrapping_rval_state_t;

mlhmmv_value_t rxval_evaluator_wrapping_rval_func(void* pvstate, variables_t* pvars) {
	rxval_evaluator_wrapping_rval_state_t* pstate = pvstate;
	rval_evaluator_t* prval_evaluator = pstate->prval_evaluator;
	mv_t val = prval_evaluator->pprocess_func(prval_evaluator->pvstate, pvars);
	return mlhmmv_value_transfer_terminal(val);
}

static void rxval_evaluator_wrapping_rval_free(rxval_evaluator_t* prxval_evaluator) {
	rxval_evaluator_wrapping_rval_state_t* pstate = prxval_evaluator->pvstate;
	pstate->prval_evaluator->pfree_func(pstate->prval_evaluator);
	free(pstate);
	free(prxval_evaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_wrapping_rval(mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr,
	int type_inferencing, int context_flags)
{
	rxval_evaluator_wrapping_rval_state_t* pstate = mlr_malloc_or_die(
		sizeof(rxval_evaluator_wrapping_rval_state_t));
	pstate->prval_evaluator = rval_evaluator_alloc_from_ast(pnode, pfmgr, type_inferencing, context_flags);

	rxval_evaluator_t* prxval_evaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	prxval_evaluator->pvstate       = pstate;
	prxval_evaluator->pprocess_func = rxval_evaluator_wrapping_rval_func;
	prxval_evaluator->pfree_func    = rxval_evaluator_wrapping_rval_free;

	return prxval_evaluator;
}

//// ================================================================
//typedef struct _rval_evaluator_from_local_variable_state_t {
//	int vardef_frame_relative_index;
//} rval_evaluator_from_local_variable_state_t;

//mv_t rval_evaluator_from_local_variable_func(void* pvstate, variables_t* pvars) {
//	rval_evaluator_from_local_variable_state_t* pstate = pvstate;
//	local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
//	mv_t val = local_stack_frame_get_non_map(pframe, pstate->vardef_frame_relative_index);
//	return mv_copy(&val);
//}

//static void rval_evaluator_from_local_variable_free(rval_evaluator_t* prxval_evaluator) {
//	rval_evaluator_from_local_variable_state_t* pstate = prxval_evaluator->pvstate;
//	free(pstate);
//	free(prxval_evaluator);
//}

//rval_evaluator_t* rval_evaluator_alloc_from_local_variable(int vardef_frame_relative_index) {
//	rval_evaluator_from_local_variable_state_t* pstate = mlr_malloc_or_die(
//		sizeof(rval_evaluator_from_local_variable_state_t));
//	rval_evaluator_t* prxval_evaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
//
//	pstate->vardef_frame_relative_index = vardef_frame_relative_index;
//	prxval_evaluator->pvstate = pstate;
//	prxval_evaluator->pprocess_func    = rval_evaluator_from_local_variable_func;
//	prxval_evaluator->pfree_func       = rval_evaluator_from_local_variable_free;
//
//	return prxval_evaluator;
//}

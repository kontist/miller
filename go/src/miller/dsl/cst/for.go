package cst

import (
	"errors"

	"miller/dsl"
	"miller/lib"
	"miller/types"
)

// ================================================================
// This is for various flavors of for-loop
// ================================================================

// ================================================================
type ForLoopKeyOnlyNode struct {
	keyVariableName    string
	mapNode            IEvaluable
	statementBlockNode *StatementBlockNode
}

func NewForLoopKeyOnlyNode(
	keyVariableName string,
	mapNode IEvaluable,
	statementBlockNode *StatementBlockNode,
) *ForLoopKeyOnlyNode {
	return &ForLoopKeyOnlyNode{
		keyVariableName,
		mapNode,
		statementBlockNode,
	}
}

// ----------------------------------------------------------------
// Sample AST:

// mlr -n put -v 'for (k in $*) { emit { k : k } }'
// DSL EXPRESSION:
// for (k in $*) { emit { k : k} }
// RAW AST:
// * StatementBlock
//     * ForLoopKeyOnly "for"
//         * LocalVariable "k"
//         * FullSrec "$*"
//         * StatementBlock
//             * EmitStatement "emit"
//                 * MapLiteral "{}"
//                     * MapLiteralKeyValuePair ":"
//                         * LocalVariable "k"
//                         * LocalVariable "k"

func BuildForLoopKeyOnlyNode(astNode *dsl.ASTNode) (*ForLoopKeyOnlyNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeForLoopKeyOnly)
	lib.InternalCodingErrorIf(len(astNode.Children) != 3)

	keyVariableASTNode := astNode.Children[0]
	mapASTNode := astNode.Children[1]
	blockASTNode := astNode.Children[2]

	lib.InternalCodingErrorIf(keyVariableASTNode.Type != dsl.NodeTypeLocalVariable)
	lib.InternalCodingErrorIf(keyVariableASTNode.Token == nil)
	keyVariableName := string(keyVariableASTNode.Token.Lit)

	// TODO: error if loop-over node isn't Mappable (inasmuch as can be
	// detected at CST-build time)
	mapNode, err := BuildEvaluableNode(mapASTNode)
	if err != nil {
		return nil, err
	}

	statementBlockNode, err := BuildStatementBlockNode(blockASTNode)
	if err != nil {
		return nil, err
	}

	return NewForLoopKeyOnlyNode(
		keyVariableName,
		mapNode,
		statementBlockNode,
	), nil
}

// ----------------------------------------------------------------
func (this *ForLoopKeyOnlyNode) Execute(state *State) error {

	mlrval := this.mapNode.Evaluate(state)
	if !mlrval.IsMap() {
		// TODO: more
		return errors.New("Miller: looped-over item is not a map.")
	}
	mapval := mlrval.GetMap()

	state.stack.PushStackFrame()
	defer state.stack.PopStackFrame()
	for pe := mapval.Head; pe != nil; pe = pe.Next {
		mapkey := types.MlrvalFromString(*pe.Key)

		// Note: The statement-block has its own push/pop for its localvars.
		// Meanwhile we need to restrict scope of the bindvar to the for-loop.
		//
		// So we have:
		//
		//   mlr put '
		//     x = 1;           <--- frame #1 main
		//     for (k in $*) {  <--- frame #2 for for-loop bindvars (right here)
		//       x = 2          <--- frame #3 for for-loop locals
		//     }
		//     x = 3;           <--- back in frame #1 main
		//   '
		//

		state.stack.BindVariable(this.keyVariableName, &mapkey)
		//state.stack.Dump()
		err := this.statementBlockNode.Execute(state)
		if err != nil {
			return err
		}
	}

	return nil
}

// ================================================================
type ForLoopKeyValueNode struct {
	keyVariableName    string
	valueVariableName  string
	mapNode            IEvaluable
	statementBlockNode *StatementBlockNode
}

func NewForLoopKeyValueNode(
	keyVariableName string,
	valueVariableName string,
	mapNode IEvaluable,
	statementBlockNode *StatementBlockNode,
) *ForLoopKeyValueNode {
	return &ForLoopKeyValueNode{
		keyVariableName,
		valueVariableName,
		mapNode,
		statementBlockNode,
	}
}

// ----------------------------------------------------------------
// Sample AST:

// mlr -n put -v 'for (k, v in $*) { emit { k : v } }'
// DSL EXPRESSION:
// for (k, v in $*) { emit { k : v } }
// RAW AST:
// * StatementBlock
//     * ForLoopKeyValue "for"
//         * LocalVariable "k"
//         * LocalVariable "v"
//         * FullSrec "$*"
//         * StatementBlock
//             * EmitStatement "emit"
//                 * MapLiteral "{}"
//                     * MapLiteralKeyValuePair ":"
//                         * LocalVariable "k"
//                         * LocalVariable "v"

func BuildForLoopKeyValueNode(astNode *dsl.ASTNode) (*ForLoopKeyValueNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeForLoopKeyValue)
	lib.InternalCodingErrorIf(len(astNode.Children) != 4)

	keyVariableASTNode := astNode.Children[0]
	valueVariableASTNode := astNode.Children[1]
	mapASTNode := astNode.Children[2]
	blockASTNode := astNode.Children[3]

	lib.InternalCodingErrorIf(keyVariableASTNode.Type != dsl.NodeTypeLocalVariable)
	lib.InternalCodingErrorIf(keyVariableASTNode.Token == nil)
	keyVariableName := string(keyVariableASTNode.Token.Lit)

	lib.InternalCodingErrorIf(valueVariableASTNode.Type != dsl.NodeTypeLocalVariable)
	lib.InternalCodingErrorIf(valueVariableASTNode.Token == nil)
	valueVariableName := string(valueVariableASTNode.Token.Lit)

	// TODO: error if loop-over node isn't Mappable (inasmuch as can be
	// detected at CST-build time)
	mapNode, err := BuildEvaluableNode(mapASTNode)
	if err != nil {
		return nil, err
	}

	statementBlockNode, err := BuildStatementBlockNode(blockASTNode)
	if err != nil {
		return nil, err
	}

	return NewForLoopKeyValueNode(
		keyVariableName,
		valueVariableName,
		mapNode,
		statementBlockNode,
	), nil
}

// ----------------------------------------------------------------
func (this *ForLoopKeyValueNode) Execute(state *State) error {

	mlrval := this.mapNode.Evaluate(state)
	if !mlrval.IsMap() {
		// TODO: more
		return errors.New("Miller: looped-over item is not a map.")
	}
	mapval := mlrval.GetMap()

	state.stack.PushStackFrame()
	defer state.stack.PopStackFrame()
	for pe := mapval.Head; pe != nil; pe = pe.Next {
		mapkey := types.MlrvalFromString(*pe.Key)

		// Note: The statement-block has its own push/pop for its localvars.
		// Meanwhile we need to restrict scope of the bindvar to the for-loop.
		//
		// So we have:
		//
		//   mlr put '
		//     x = 1;           <--- frame #1 main
		//     for (k in $*) {  <--- frame #2 for for-loop bindvars (right here)
		//       x = 2          <--- frame #3 for for-loop locals
		//     }
		//     x = 3;           <--- back in frame #1 main
		//   '
		//

		state.stack.BindVariable(this.keyVariableName, &mapkey)
		state.stack.BindVariable(this.valueVariableName, pe.Value)
		//state.stack.Dump()
		err := this.statementBlockNode.Execute(state)
		if err != nil {
			return err
		}
	}

	return nil
}

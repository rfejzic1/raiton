package parser

// Implements Expression
type Scope struct {
	typeDefinitions []TypeDefinition
	definitions     []Definition
	expressions     []Expression
}

type TypeDefinition struct {
	identifier     TypeIdentifier
	typeExpression TypeExpression
}

type Definition struct {
	typeExpression TypeExpression // if not defined explicitly, inferred from expression
	identifier     Identifier
	parameters     []Identifier
	expression     Expression
}

type Expression interface{}

type TypeExpression interface{}

// *** Type Expressions ***

type TypeIdentifier string // e.g. string, number, list, person, etc.

type FunctionType struct {
	parameterType TypeExpression
	returnType    TypeExpression
} // e.g. number -> (number -> number); argumentType -> returnType

type RecordType struct {
	fields map[Identifier]TypeExpression
} // e.g. { name: string, age: number }

// TODO:
// type EnumType struct{} // e.g. Red | Green | Blue | RGB { r: number, g: number, b: number}

// *** Expressions ***

type Identifier string

type Invocation struct {
	invokee   Expression
	arguments []Expression
}

type LambdaLiteral struct {
	parameters []Identifier
	expression Expression
}

type RecordLiteral struct {
	fields map[Identifier]Expression
}

type ListLiteral struct {
	elements []Expression
}

type NumberLiteral string

type StringLiteral string

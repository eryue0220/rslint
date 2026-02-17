package no_dupe_keys

import (
	"testing"

	"github.com/web-infra-dev/rslint/internal/plugins/typescript/rules/fixtures"
	"github.com/web-infra-dev/rslint/internal/rule_tester"
)

func TestNoDupeKeysRule(t *testing.T) {
	rule_tester.RunRuleTester(
		fixtures.GetRootDir(),
		"tsconfig.json",
		t,
		&NoDupeKeysRule,
		// Valid cases - ported from ESLint
		[]rule_tester.ValidTestCase{
			{Code: "var foo = { __proto__: 1, two: 2};"},
			{Code: "var x = { foo: 1, bar: 2 };"},
			{Code: "var x = { '': 1, bar: 2 };"},
			{Code: "var x = { '': 1, ' ': 2 };"},
			{Code: "var x = { '': 1, [null]: 2 };"},
			{Code: "var x = { '': 1, [a]: 2 };"},
			{Code: "var x = { a: b, [a]: b };"},
			{Code: "var x = { a: b, ...c }"},
			{Code: "var x = { [a]: 1, [a]: 2 };"},
			{Code: "+{ get a() { }, set a(b) { } };"},
			{Code: "var x = { get a() {}, set a (value) {} };"},
			{Code: "var x = { a: 1, b: { a: 2 } };"},
			{Code: "var x = ({ null: 1, [/(?<zero>0)/]: 2 })"},
			{Code: "var {a, a} = obj"},
			{Code: "var x = { 1_0: 1, 1: 2 };"},
			{Code: "var x = { __proto__: null, ['__proto__']: null };"},
			{Code: "var x = { ['__proto__']: null, __proto__: null };"},
			{Code: "var x = { '__proto__': null, ['__proto__']: null };"},
			{Code: "var x = { ['__proto__']: null, '__proto__': null };"},
			{Code: "var x = { __proto__: null, __proto__ };"},
			{Code: "var x = { __proto__, __proto__: null };"},
			{Code: "var x = { __proto__: null, __proto__() {} };"},
			{Code: "var x = { __proto__() {}, __proto__: null };"},
			{Code: "var x = { __proto__: null, get __proto__() {} };"},
			{Code: "var x = { get __proto__() {}, __proto__: null };"},
			{Code: "var x = { __proto__: null, set __proto__(value) {} };"},
			{Code: "var x = { set __proto__(value) {}, __proto__: null };"},

			// Syntax Error: Octal literals are not allowed
			// {Code: "var x = { 012: 1, 12: 2 };"},
		},
		// Invalid cases - ported from ESLint
		[]rule_tester.InvalidTestCase{
			{
				Code: "var x = { a: b, ['a']: b };",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 1, Column: 17},
				},
			},
			{
				Code: "var x = { y: 1, y: 2 };",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 1, Column: 17},
				},
			},
			{
				Code: "var x = { '': 1, '': 2 };",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 1, Column: 18},
				},
			},
			{
				Code: "var x = { '': 1, [``]: 2 };",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 1, Column: 18},
				},
			},
			{
				Code: "var foo = { 0x1: 1, 1: 2};",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 1, Column: 21},
				},
			},
			{
				Code: "var x = { 0b1: 1, 1: 2 };",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 1, Column: 19},
				},
			},
			{
				Code: "var x = { 0o1: 1, 1: 2 };",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 1, Column: 19},
				},
			},
			{
				Code: "var x = { 1n: 1, 1: 2 };",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 1, Column: 18},
				},
			},
			{
				Code: "var x = { 1_0: 1, 10: 2 };",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 1, Column: 19},
				},
			},
			{
				Code: "var x = { \"z\": 1, z: 2 };",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 1, Column: 19},
				},
			},
			{
				Code: "var foo = {\n  bar: 1,\n  bar: 1,\n}",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 3, Column: 3},
				},
			},
			{
				Code: "var x = { a: 1, get a() {} };",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 1, Column: 17},
				},
			},
			{
				Code: "var x = { a: 1, set a(value) {} };",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 1, Column: 17},
				},
			},
			{
				Code: "var x = { a: 1, b: { a: 2 }, get b() {} };",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 1, Column: 30},
				},
			},
			{
				Code: "var x = ({ '/(?<zero>0)/': 1, [/(?<zero>0)/]: 2 })",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 1, Column: 31},
				},
			},
			{
				Code: "var x = { ['__proto__']: null, ['__proto__']: null };",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 1, Column: 32},
				},
			},
			{
				Code: "var x = { ['__proto__']: null, __proto__ };",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 1, Column: 32},
				},
			},
			{
				Code: "var x = { ['__proto__']: null, __proto__() {} };",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 1, Column: 32},
				},
			},
			{
				Code: "var x = { ['__proto__']: null, get __proto__() {} };",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 1, Column: 32},
				},
			},
			{
				Code: "var x = { ['__proto__']: null, set __proto__(value) {} };",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 1, Column: 32},
				},
			},
			{
				Code: "var x = { __proto__: null, a: 5, a: 6 };",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpected", Line: 1, Column: 34},
				},
			},

			// Syntax Error: Octal literals are not allowed
			// {
			// 	Code: "var x = { 012: 1, 10: 2 };",
			// 	Errors: []rule_tester.InvalidTestCaseError{
			// 		{
			// 			MessageId: "unexpected",
			// 			Line:      1,
			// 			Column:    1,
			// 		},
			// 	},
			// },
		},
	)
}

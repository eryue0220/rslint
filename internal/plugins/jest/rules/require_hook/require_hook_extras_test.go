// TestRequireHookExtras locks in branches and edge shapes that the upstream
// test suite doesn't exercise. Each case carries an inline comment pointing at
// the specific branch / Dimension 4 row / tsgo AST quirk it covers, so future
// refactors can't silently regress them without breaking a named lock-in.

package require_hook_test

import (
	"testing"

	"github.com/web-infra-dev/rslint/internal/plugins/jest/fixtures"
	"github.com/web-infra-dev/rslint/internal/plugins/jest/rules/require_hook"
	"github.com/web-infra-dev/rslint/internal/rule_tester"
)

func TestRequireHookExtras(t *testing.T) {
	rule_tester.RunRuleTester(
		fixtures.GetRootDir(),
		"tsconfig.json",
		t,
		&require_hook.RequireHookRule,
		[]rule_tester.ValidTestCase{
			// ---- Dimension 4: parenthesized call is not flagged ----
			// Upstream shouldBeInHook only unwraps ExpressionStatement once; a
			// ParenthesizedExpression hits the default arm and is ignored.
			{Code: `(setup());`},

			// ---- Dimension 4: optional-chain call is still a CallExpression ----
			// N/A for valid — covered as invalid below.

			// ---- Dimension 4: element-access jest call is allowed ----
			{Code: `jest['mock']('../api');`},

			// ---- Dimension 4: describe concise arrow body is skipped ----
			// getFunctionBodyBlock requires a Block; expression-body describe is ignored.
			{Code: `describe('a test', () => setup());`},

			// ---- Dimension 4: function declaration at top level is allowed ----
			// Locks in upstream shouldBeInHook() default arm for FunctionDeclaration.
			{Code: `function helper() { setup(); }
test('x', () => {});`},

			// ---- Dimension 4: assignment expression is allowed ----
			// Locks in upstream shouldBeInHook() default arm for AssignmentExpression.
			{Code: `let x;
x = 1;
test('x', () => {});`},

			// ---- Dimension 4: new expression alone is allowed ----
			{Code: `new Helper();
test('x', () => {});`},

			// ---- Dimension 4: empty describe callback body ----
			{Code: `describe('empty', () => {});`},

			// ---- Dimension 4: describe.only still inspected, hooks allowed ----
			{Code: `describe.only('suite', () => {
  beforeEach(() => setup());
});`},

			// ---- Dimension 4: var without initializer is allowed ----
			{Code: `var x;
test('x', () => {});`},

			// ---- Dimension 4: type alias / interface are allowed ----
			{Code: `type T = number;
interface I { x: number }
test('x', () => {});`},

			// Locks in upstream isJestFnCall() via ParseJestFnCall for custom jest.* APIs.
			{Code: `jest.anythingCustom();`},

			// Locks in upstream shouldBeInHook() const arm — require() in const is fine.
			{Code: `const utils = require('./utils');
test('x', () => {});`},

			// Locks in upstream CallExpression listener: describe with < 2 args is skipped.
			{Code: `describe('title only');`},

			// Locks in upstream CallExpression listener: non-function second arg is skipped.
			{Code: `describe('title', 'not a function');`},

			// ---- Real-user: issue#934 overrides / non-test file noise ----
			// Top-level side effects in non-test files are still flagged by the rule;
			// users restrict via overrides. Lock the aggressive top-level behavior.
			// (Covered by invalid `setup();` upstream case.)

			// ---- Real-user: issue#1386 chained tester after new ----
			// Valid workaround: wrap in beforeEach / allowedFunctionCalls.
			{
				Code: `enableAutoDestroy(afterEach);
test('x', () => {});`,
				Options: []interface{}{
					map[string]interface{}{
						"allowedFunctionCalls": []interface{}{"enableAutoDestroy"},
					},
				},
			},

			// Locks in upstream allowedFunctionCalls exact-name match for dotted callee.
			{
				Code: `helper.setup();`,
				Options: []interface{}{
					map[string]interface{}{
						"allowedFunctionCalls": []interface{}{"helper.setup"},
					},
				},
			},
		},
		[]rule_tester.InvalidTestCase{
			// ---- Dimension 4: optional-chain call ----
			{
				Code: `setup?.();`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Line: 1, Column: 1},
				},
			},

			// ---- Dimension 4: describe function expression body ----
			{
				Code: `describe('suite', function () {
  setup();
});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Line: 2, Column: 3},
				},
			},

			// ---- Dimension 4: var with initializer ----
			// Locks in upstream shouldBeInHook() VariableDeclaration non-const arm.
			{
				Code: `var value = 1;`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Line: 1, Column: 1},
				},
			},

			// ---- Dimension 4: parenthesized null initializer is NOT treated as null ----
			// Upstream isNullOrUndefined does not unwrap parentheses.
			{
				Code: `let x = (null);`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Line: 1, Column: 1},
				},
			},

			// ---- Dimension 4: member call chain ----
			{
				Code: `foo.bar();`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Line: 1, Column: 1},
				},
			},

			// ---- Dimension 4: nested describe.only ----
			{
				Code: `describe.only('suite', () => {
  setup();
});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Line: 2, Column: 3},
				},
			},

			// Locks in upstream isJestFnCall() / allowedFunctionCalls miss for dotted names.
			{
				Code: `helper.setup();`,
				Options: []interface{}{
					map[string]interface{}{
						"allowedFunctionCalls": []interface{}{"setup"},
					},
				},
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Line: 1, Column: 1},
				},
			},

			// ---- Real-user: issue#1386 class method chain triggers require-hook ----
			{
				Code: `new NodeExtensionTester()
  .shouldMatch()
  .runTests();`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Line: 1, Column: 1},
				},
			},

			// Calls inside test bodies must not be reported (only Program / describe bodies).
			// Locks in that the CallExpression listener only expands describe callbacks.
			{
				Code: `test('x', () => {
  setup();
});
setup();`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Line: 4, Column: 1},
				},
			},
		},
	)
}

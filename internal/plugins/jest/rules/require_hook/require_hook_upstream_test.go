// TestRequireHookUpstream migrates the full valid/invalid suite from upstream
// src/rules/__tests__/require-hook.test.ts 1:1. Position assertions cover
// line/column for every invalid case. rslint-specific lock-in cases live in the
// require_hook_extras_test.go file(s).

package require_hook_test

import (
	"testing"

	"github.com/web-infra-dev/rslint/internal/plugins/jest/fixtures"
	"github.com/web-infra-dev/rslint/internal/plugins/jest/rules/require_hook"
	"github.com/web-infra-dev/rslint/internal/rule_tester"
)

func TestRequireHookUpstream(t *testing.T) {
	rule_tester.RunRuleTester(
		fixtures.GetRootDir(),
		"tsconfig.json",
		t,
		&require_hook.RequireHookRule,
		[]rule_tester.ValidTestCase{
			// ---- default ----
			{Code: "describe()"},
			{Code: `describe("just a title")`},
			{Code: `describe('a test', () =>
  test('something', () => {
    expect(true).toBe(true);
  }));`},
			{Code: `test('it', () => {
  //
});`},
			{Code: `const { myFn } = require('../functions');

test('myFn', () => {
  expect(myFn()).toBe(1);
});`},
			{Code: `import { myFn } from '../functions';

test('myFn', () => {
  expect(myFn()).toBe(1);
});`},
			{Code: `class MockLogger {
  log() {}
}

test('myFn', () => {
  expect(myFn()).toBe(1);
});`},
			{Code: `const { myFn } = require('../functions');

describe('myFn', () => {
  it('returns one', () => {
    expect(myFn()).toBe(1);
  });
});`},
			{Code: `describe('some tests', () => {
  it('is true', () => {
    expect(true).toBe(true);
  });
});`},
			{Code: `describe('some tests', () => {
  it('is true', () => {
    expect(true).toBe(true);
  });

  describe('more tests', () => {
    it('is false', () => {
      expect(true).toBe(false);
    });
  });
});`},
			{Code: `describe('some tests', () => {
  let consoleLogSpy;

  beforeEach(() => {
    consoleLogSpy = jest.spyOn(console, 'log');
  });

  it('prints a message', () => {
    printMessage('hello world');

    expect(consoleLogSpy).toHaveBeenCalledWith('hello world');
  });
});`},
			{Code: `let consoleErrorSpy = null;

beforeEach(() => {
  consoleErrorSpy = jest.spyOn(console, 'error');
});`},
			{Code: `let consoleErrorSpy = undefined;

beforeEach(() => {
  consoleErrorSpy = jest.spyOn(console, 'error');
});`},
			{Code: `describe('some tests', () => {
  beforeEach(() => {
    setup();
  });
});`},
			{Code: `beforeEach(() => {
  initializeCityDatabase();
});

afterEach(() => {
  clearCityDatabase();
});

test('city database has Vienna', () => {
  expect(isCity('Vienna')).toBeTruthy();
});

test('city database has San Juan', () => {
  expect(isCity('San Juan')).toBeTruthy();
});`},
			{Code: `describe('cities', () => {
  beforeEach(() => {
    initializeCityDatabase();
  });

  test('city database has Vienna', () => {
    expect(isCity('Vienna')).toBeTruthy();
  });

  test('city database has San Juan', () => {
    expect(isCity('San Juan')).toBeTruthy();
  });

  afterEach(() => {
    clearCityDatabase();
  });
});`},
			{
				Code: `enableAutoDestroy(afterEach);

describe('some tests', () => {
  it('is false', () => {
    expect(true).toBe(true);
  });
});`,
				Options: []interface{}{
					map[string]interface{}{
						"allowedFunctionCalls": []interface{}{"enableAutoDestroy"},
					},
				},
			},
			// ---- typescript edition ----
			{Code: `import { myFn } from '../functions';

// todo: https://github.com/DefinitelyTyped/DefinitelyTyped/pull/56545
declare module 'eslint' {
  namespace ESLint {
    interface LintResult {
      fatalErrorCount: number;
    }
  }
}

test('myFn', () => {
  expect(myFn()).toBe(1);
});`},
		},
		[]rule_tester.InvalidTestCase{
			// ---- default ----
			{
				Code: "setup();",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Message: "This should be done within a hook", Line: 1, Column: 1},
				},
			},
			{
				Code: `describe('some tests', () => {
  setup();
});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Line: 2, Column: 3},
				},
			},
			{
				Code: `let { setup } = require('./test-utils');

describe('some tests', () => {
  setup();
});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Line: 1, Column: 1},
					{MessageId: "useHook", Line: 4, Column: 3},
				},
			},
			{
				Code: `describe('some tests', () => {
  setup();

  it('is true', () => {
    expect(true).toBe(true);
  });

  describe('more tests', () => {
    setup();

    it('is false', () => {
      expect(true).toBe(false);
    });
  });
});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Line: 2, Column: 3},
					{MessageId: "useHook", Line: 9, Column: 5},
				},
			},
			{
				Code: `let consoleErrorSpy = jest.spyOn(console, 'error');

describe('when loading cities from the api', () => {
  let consoleWarnSpy = jest.spyOn(console, 'warn');
});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Line: 1, Column: 1},
					{MessageId: "useHook", Line: 4, Column: 3},
				},
			},
			{
				Code: `let consoleErrorSpy = null;

describe('when loading cities from the api', () => {
  let consoleWarnSpy = jest.spyOn(console, 'warn');
});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Line: 4, Column: 3},
				},
			},
			{
				Code: "let value = 1",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Line: 1, Column: 1},
				},
			},
			{
				Code: "let consoleErrorSpy, consoleWarnSpy = jest.spyOn(console, 'error');",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Line: 1, Column: 1},
				},
			},
			{
				Code: "let consoleErrorSpy = jest.spyOn(console, 'error'), consoleWarnSpy;",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Line: 1, Column: 1},
				},
			},
			{
				Code: `import { database, isCity } from '../database';
import { loadCities } from '../api';

jest.mock('../api');

const initializeCityDatabase = () => {
  database.addCity('Vienna');
  database.addCity('San Juan');
  database.addCity('Wellington');
};

const clearCityDatabase = () => {
  database.clear();
};

initializeCityDatabase();

test('that persists cities', () => {
  expect(database.cities.length).toHaveLength(3);
});

test('city database has Vienna', () => {
  expect(isCity('Vienna')).toBeTruthy();
});

test('city database has San Juan', () => {
  expect(isCity('San Juan')).toBeTruthy();
});

describe('when loading cities from the api', () => {
  let consoleWarnSpy = jest.spyOn(console, 'warn');

  loadCities.mockResolvedValue(['Wellington', 'London']);

  it('does not duplicate cities', async () => {
    await database.loadCities();

    expect(database.cities).toHaveLength(4);
  });

  it('logs any duplicates', async () => {
    await database.loadCities();

    expect(consoleWarnSpy).toHaveBeenCalledWith(
      'Ignored duplicate cities: Wellington',
    );
  });
});

clearCityDatabase();`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Line: 16, Column: 1},
					{MessageId: "useHook", Line: 31, Column: 3},
					{MessageId: "useHook", Line: 33, Column: 3},
					{MessageId: "useHook", Line: 50, Column: 1},
				},
			},
			{
				Code: `enableAutoDestroy(afterEach);

describe('some tests', () => {
  it('is false', () => {
    expect(true).toBe(true);
  });
});`,
				Options: []interface{}{
					map[string]interface{}{
						"allowedFunctionCalls": []interface{}{"someOtherName"},
					},
				},
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Line: 1, Column: 1},
				},
			},
			// ---- typescript edition ----
			{
				Code: `import { setup } from '../test-utils';

// todo: https://github.com/DefinitelyTyped/DefinitelyTyped/pull/56545
declare module 'eslint' {
  namespace ESLint {
    interface LintResult {
      fatalErrorCount: number;
    }
  }
}

describe('some tests', () => {
  setup();
});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "useHook", Line: 13, Column: 3},
				},
			},
		},
	)
}

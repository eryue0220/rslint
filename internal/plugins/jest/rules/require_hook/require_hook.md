# require-hook

## Rule Details

Require setup and teardown code to be within a Jest lifecycle hook (`beforeEach`, `afterEach`, `beforeAll`, `afterAll`).

Jest runs all `describe` handlers before any tests, so top-level or `describe`-body side effects (non-`const` initializers, non-Jest function calls) should live in hooks instead.

Examples of **incorrect** code for this rule:

```javascript
setup();

describe('suite', () => {
  setup();

  it('works', () => {
    expect(true).toBe(true);
  });
});

let consoleErrorSpy = jest.spyOn(console, 'error');
```

Examples of **correct** code for this rule:

```javascript
import { myFn } from '../functions';

const helper = () => myFn();

beforeEach(() => {
  setup();
});

describe('suite', () => {
  let consoleLogSpy;

  beforeEach(() => {
    consoleLogSpy = jest.spyOn(console, 'log');
  });

  it('works', () => {
    expect(myFn()).toBe(1);
  });
});
```

## Options

- First argument (optional): object with `allowedFunctionCalls`
  - `allowedFunctionCalls`: array of callee names that are allowed outside hooks (e.g. `enableAutoDestroy`).

```json
{ "jest/require-hook": ["error", { "allowedFunctionCalls": ["enableAutoDestroy"] }] }
```

```javascript
enableAutoDestroy(afterEach);

describe('suite', () => {
  it('works', () => {
    expect(true).toBe(true);
  });
});
```

## Original Documentation

- [jest/require-hook](https://github.com/jest-community/eslint-plugin-jest/blob/main/docs/rules/require-hook.md)

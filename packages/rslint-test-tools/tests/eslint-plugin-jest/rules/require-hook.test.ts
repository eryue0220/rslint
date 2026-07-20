import { RuleTester } from '../rule-tester';

const ruleTester = new RuleTester();

ruleTester.run('require-hook', {} as never, {
  valid: [
    { code: 'describe()' },
    { code: 'describe("just a title")' },
    {
      code: `
      describe('a test', () =>
        test('something', () => {
          expect(true).toBe(true);
        }));
    `,
    },
    {
      code: `
      test('it', () => {
        //
      });
    `,
    },
    {
      code: `
      const { myFn } = require('../functions');

      test('myFn', () => {
        expect(myFn()).toBe(1);
      });
    `,
    },
    {
      code: `
        import { myFn } from '../functions';

        test('myFn', () => {
          expect(myFn()).toBe(1);
        });
      `,
    },
    {
      code: `
      class MockLogger {
        log() {}
      }

      test('myFn', () => {
        expect(myFn()).toBe(1);
      });
    `,
    },
    {
      code: `
      const { myFn } = require('../functions');

      describe('myFn', () => {
        it('returns one', () => {
          expect(myFn()).toBe(1);
        });
      });
    `,
    },
    {
      code: `
      describe('some tests', () => {
        it('is true', () => {
          expect(true).toBe(true);
        });
      });
    `,
    },
    {
      code: `
      describe('some tests', () => {
        it('is true', () => {
          expect(true).toBe(true);
        });

        describe('more tests', () => {
          it('is false', () => {
            expect(true).toBe(false);
          });
        });
      });
    `,
    },
    {
      code: `
      describe('some tests', () => {
        let consoleLogSpy;

        beforeEach(() => {
          consoleLogSpy = jest.spyOn(console, 'log');
        });

        it('prints a message', () => {
          printMessage('hello world');

          expect(consoleLogSpy).toHaveBeenCalledWith('hello world');
        });
      });
    `,
    },
    {
      code: `
      let consoleErrorSpy = null;

      beforeEach(() => {
        consoleErrorSpy = jest.spyOn(console, 'error');
      });
    `,
    },
    {
      code: `
      let consoleErrorSpy = undefined;

      beforeEach(() => {
        consoleErrorSpy = jest.spyOn(console, 'error');
      });
    `,
    },
    {
      code: `
      describe('some tests', () => {
        beforeEach(() => {
          setup();
        });
      });
    `,
    },
    {
      code: `
      beforeEach(() => {
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
      });
    `,
    },
    {
      code: `
      describe('cities', () => {
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
      });
    `,
    },
    {
      code: `
        enableAutoDestroy(afterEach);

        describe('some tests', () => {
          it('is false', () => {
            expect(true).toBe(true);
          });
        });
      `,
      options: [{ allowedFunctionCalls: ['enableAutoDestroy'] }],
    },
    {
      code: `
      import { myFn } from '../functions';

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
      });
    `,
    },
  ],
  invalid: [
    {
      code: 'setup();',
      errors: [
        {
          messageId: 'useHook',
          line: 1,
          column: 1,
        },
      ],
    },
    {
      code: `
        describe('some tests', () => {
          setup();
        });
      `,
      errors: [
        {
          messageId: 'useHook',
          line: 2,
          column: 3,
        },
      ],
    },
    {
      code: `
        let { setup } = require('./test-utils');

        describe('some tests', () => {
          setup();
        });
      `,
      errors: [
        {
          messageId: 'useHook',
          line: 1,
          column: 1,
        },
        {
          messageId: 'useHook',
          line: 4,
          column: 3,
        },
      ],
    },
    {
      code: `
        describe('some tests', () => {
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
        });
      `,
      errors: [
        {
          messageId: 'useHook',
          line: 2,
          column: 3,
        },
        {
          messageId: 'useHook',
          line: 9,
          column: 5,
        },
      ],
    },
    {
      code: `
        let consoleErrorSpy = jest.spyOn(console, 'error');

        describe('when loading cities from the api', () => {
          let consoleWarnSpy = jest.spyOn(console, 'warn');
        });
      `,
      errors: [
        {
          messageId: 'useHook',
          line: 1,
          column: 1,
        },
        {
          messageId: 'useHook',
          line: 4,
          column: 3,
        },
      ],
    },
    {
      code: `
        let consoleErrorSpy = null;

        describe('when loading cities from the api', () => {
          let consoleWarnSpy = jest.spyOn(console, 'warn');
        });
      `,
      errors: [
        {
          messageId: 'useHook',
          line: 4,
          column: 3,
        },
      ],
    },
    {
      code: 'let value = 1',
      errors: [
        {
          messageId: 'useHook',
          line: 1,
          column: 1,
        },
      ],
    },
    {
      code: "let consoleErrorSpy, consoleWarnSpy = jest.spyOn(console, 'error');",
      errors: [
        {
          messageId: 'useHook',
          line: 1,
          column: 1,
        },
      ],
    },
    {
      code: "let consoleErrorSpy = jest.spyOn(console, 'error'), consoleWarnSpy;",
      errors: [
        {
          messageId: 'useHook',
          line: 1,
          column: 1,
        },
      ],
    },
    {
      code: `
        import { database, isCity } from '../database';
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

        clearCityDatabase();
      `,
      errors: [
        {
          messageId: 'useHook',
          line: 16,
          column: 1,
        },
        {
          messageId: 'useHook',
          line: 31,
          column: 3,
        },
        {
          messageId: 'useHook',
          line: 33,
          column: 3,
        },
        {
          messageId: 'useHook',
          line: 50,
          column: 1,
        },
      ],
    },
    {
      code: `
        enableAutoDestroy(afterEach);

        describe('some tests', () => {
          it('is false', () => {
            expect(true).toBe(true);
          });
        });
      `,
      options: [{ allowedFunctionCalls: ['someOtherName'] }],
      errors: [
        {
          messageId: 'useHook',
          line: 1,
          column: 1,
        },
      ],
    },
    {
      code: `
        import { setup } from '../test-utils';

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
        });
      `,
      errors: [
        {
          messageId: 'useHook',
          line: 13,
          column: 3,
        },
      ],
    },
  ],
});

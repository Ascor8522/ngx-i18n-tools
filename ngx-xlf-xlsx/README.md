# ngx-xlf-xlsx

[![NPM Version](https://img.shields.io/npm/v/ngx-xlf-xlsx)](https://www.npmjs.com/package/ngx-xlf-xlsx)
[![GitHub](https://img.shields.io/badge/GitHub-%23121011.svg?&logo=github)](https://github.com/Ascor8522/ngx-i18n-tools)

**A simple CLI to convert back and forth between XLIFF and XLSX files in Angular projects.**

## Workflow / How to Use

In your Angular project's root directory, run:

1. Extract the source strings from the Angular app:

   ```bash
   ng extract-i18n --output-path src/locale
   ```

2. Run the CLI:

   ```bash
   npx ngx-xlf-xlsx@latest
   ```

   No (prior) installation of the tool is required.

   For other package managers:

   - pnpm: `pnpx ngx-xlf-xlsx@latest`
   - yarn: `yarn dlx ngx-xlf-xlsx@latest`

3. Translate the strings in the Excel file in `translations.xlsx`.

   The file is automatically created in the root of the project if needed.

4. Run the CLI again:

   ```bash
   npx ngx-xlf-xlsx@latest
   ```

## Requirements, Assumptions, and Precautions

- `translations.xlsx` will be used for translations.
   If the file already exists and contains other kinds of data, it will be overwritten.
- The CLI will remove any obsolete translations from the Excel file.
- In your

## Error Reference

### `json: cannot unmarshal object into Go struct field .projects.i18n.locales of type main.Path`

This error occurs when the `angular.json` config file does not have the expected shape for the i18n configuration.

the locales must be specified as a key-value object, where the key is the locale name and the value is the path to the XLF file.

Example:

```json
{
 "projects": {
  "i18n": {
   "locales": {
    "fr": "src/locale/messages.fr.xlf",
    "nl": "src/locale/messages.nl.xlf",
    "de": "src/locale/messages.de.xlf"
   }
  }
 }
}
```

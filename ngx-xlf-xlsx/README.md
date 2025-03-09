# ngx-xlf-xlsx

[![NPM Version](https://img.shields.io/npm/v/ngx-xlf-xlsx)](https://www.npmjs.com/package/ngx-xlf-xlsx)
[![GitHub](https://img.shields.io/badge/GitHub-%23121011.svg?&logo=github)](https://github.com/Ascor8522/ngx-i18n-tools)

A simple CLI to convert back and forth between XLIFF (`.xlf`) and Excel (`.xlsx`) files in Angular projects.

## Installation

Installation, either local or global, is possible, though not recommended.

Instead, call the tool using `npx` and ask for the latest version, as shown in the next section.

```bash
npm install ngx-xlf-xlsx
```

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

- The Angular project is using `@angular/localize` to manage internationalization.
- The project is configured with the different locales in `angular.json` file, as key-value pairs,
where the key is the locale name and the value is the path to the XLF file.
- `translations.xlsx` is used for translations.
If the file already exists and contains other kinds of data, it will be overwritten and the data will be LOST.
- The CLI will remove any obsolete translations from the Excel file.
- The CLI will overwrite the content of the non-source XLF files.
Make sure to back them up if you want to keep them.
If you want to keep them, make a backup of the file before running the CLI.
- Source strings cannot contain patterns like `${{variable}}`,
as they are considered as placeholders and will be attempted to be deserialized as such.
- Strings with multiple placeholders have to have names specified for each placeholder.

## Migration

If you are migrating from Walter's tool™️, check out the [migration guide](./MIGRATION.md).

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

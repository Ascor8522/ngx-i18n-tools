# Migration Guide

This tool was meant to replace an internal tool made by Walter.

## Feature Comparison

|                                      | Walter's Tool™️  | ngx-xlf-xlsx |
|--------------------------------------|-----------------|--------------|
| Works out of the box                 | <center>❌<sup>1 | <center>✔️    |
| Easy to use                          | <center>❌<sup>2 | <center>✔️    |
| Easy to add more languages           | <center>❌<sup>3 | <center>✔️    |
| Updates source strings in Excel      | <center>❌       | <center>✔️    |
| Removes unused strings from Excel    | <center>❌       | <center>✔️    |
| Handles placeholders correctly       | <center>❌<sup>4 | <center>✔️    |
| Cleans up translations from business | <center>❌<sup>5 | <center>✔️    |
| Made by Walter                       | <center>✔️       | <center>❌    |

1.
   - requires having the JRE installed and in PATH
   - requires having Apache Ant installed and in PATH
   - requires having the necessary configuration files committed to the repo
   - requires having the necessary libraries in the right place and be committed to the repo
2.
   - involves copying files around
   - involves running commands in the right order
   - generates lots of temporary files
   - requires Excel file to exist, and does not create it automatically
3.
   - requires updating the Ant build config file and the various shell scripts
   - does not follow config described in `angular.json`
4.
   - does not display placeholders correctly in Excel file
   - does not ensure number of placeholders is consistent between source and target strings
5.
   - does not remove trailing whitespace
   - does not remove multiple spaces in a row
   - does not replace unicode characters with their ASCII equivalent (TODO)

## Migration Steps

1. Adapt your `angular.json` config file to use a key-value object for the i18n locales.
   The keys are the locales, and the values are the paths to the XLF files, relative to the project root.

   from:

   ```json
   {
    "projects": {
     "i18n": {
      "sourceLocale": "en",
      "locales": {
       "fr": {
        "translation": "src/locale/messages.fr.xlf"
       },
       "nl": {
        "translation": "src/locale/messages.nl.xlf"
       },
       "de": {
        "translation": "src/locale/messages.de.xlf"
       }
      }
     }
    }
   }

   ```

   to:

   ```json
   {
    "projects": {
     "i18n": {
      "sourceLocale": "en",
      "locales": {
       "fr": "src/locale/messages.fr.xlf",
       "nl": "src/locale/messages.nl.xlf",
       "de": "src/locale/messages.de.xlf"
      }
     }
    }
   }
   ```

2. Move the `translations/data/messages.xlsx` file to `translations.xlsx` in the root of the project.

   ```bash
   mv translations/data/messages.xlsx translations.xlsx
   ```

3. Remove the `translations` directory.

   ```bash
   rm -rf translations
   ```

4. Remove `build.xml` file in the root of the project.

   ```bash
   rm build.xml
   ```

5. In your source files,
   1. Use actual placeholders instead of using `.replace("$placeholder", value)` on a translation.

      from:

      ```ts
      $localize`Hello $name`.replace("$name", name);
      ```

      to:

      ```ts
      $localize`Hello ${name}:name:`;
      ```

   2. Make sure your placeholders are named.
   You can name a placeholder by using the `$localize``Hello ${name}:name:``.` syntax.
   Notice the `:name:` part after the placeholder.

      from:

      ```ts
      $localize`Hello ${name}`
      ```

      to:

      ```ts
      $localize`Hello ${name}:name:`
      ```

   3. Make sure you don't have any `${{placeholder}}` patterns in your source strings.
   If you do, replace them with a different pattern.

      from:

      ```html
      <p>Price: ${{ amount | number }}</p>
      ```

      to:

      ```html
      <p>Price: {{ amount | currency }}</p>
      ```

6. In `translations.xlsx`,
   1. Replace all `|` characters with line actual breaks.
   You can insert line breaks in Excel by pressing <kbd>Alt+Enter</kbd>.

      from:

      ```text
      Bonjour ${{name}},|Bienvenue sur notre site web!
      ```

      to:

      ```text
      Bonjour ${{name}},
      Bienvenue sur notre site web!
      ```

   2. Adapt your placeholders to use the new `${{PLACEHOLDER_NAME}}` syntax.

      from:

      ```text
      Bonjour $name
      ```

      to:

      ```text
      Bonjour ${{name}}
      ```

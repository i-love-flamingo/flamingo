# OM3 Coding Conventions

## Project Structure

The project structure is meant to be component oriented. That means that we have
separate folders for the different components we have and in these folders
should be everything the component needs.

For example we have a button component like this:

+ ../frontend/src/component/button
  + button.md
  + button.pug
  + button.sass
  + button.vue
  + (button.mock.json)
  + (button.spec.js)

To give us nice search hits the files should all be named / prefixed with the
component name - if yoou search for button, you will get all the files for the button
component.

# Change Log - Pglet Server

## [0.5.7](https://github.com/pglet/pglet/releases/tag/v0.5.7)

* `Stack` automatically scrolls to bottom if `autoscroll` property set to `true`.
* Set `page.UserAuthProvider` to a used authentication method (`github`, `google` or `azure`).
* `page.win_width` and `page.win_height` properties renamed to `page.winwidth` and `page.winheight`.
* When host is connected to a `page` its contents and properties are cleaned unless `update: true` is passed. No need to call `page.clean()` on the client anymore.
* Focusing input controls - allows setting focus on a control when added to a page or page loaded:
  * `Button.focused`
  * `Checkbox.focused`
  * `ChoiceGroup.focused`
  * `DatePicker.focused`
  * `Dropdown.focused`
  * `SearchBox.focused`
  * `Slider.focused`
  * `SpinButton.focused`
  * `Textbox.focused`
  * `Toggle.focused`
* `focus` and `blur` events for the following input controls:
  * `Button`
  * `ChoiceGroup`
  * `DatePicker`
  * `Dropdown`
  * `SearchBox`
  * `Slider`
  * `SpinButton`
  * `Textbox`
  * `Toggle`
* New `page` properties:
  * `userAuthProvider`
* New `page` events:
  * `connect` - web client connected
  * `disconnect` - web client disconnected
* New `IFrame` properties:
  * `borderWidth`
  * `borderColor`
  * `borderStyle`
  * `borderRadius`
* New `Stack` properties:
  * `autoscroll`
  * `borderWidth`
  * `borderColor`
  * `borderStyle`
* New `Image` properties:
  * `fit` = `none`, `contain`, `cover`, `center`, `centerContain`, `centerCover`
  * `borderWidth`
  * `borderColor`
  * `borderStyle`
  * `borderRadius`
* New `Persona` control:
  * `imageUrl`
  * `imageAlt`
  * `initialsColor`
  * `initialsTextColor`
  * `text`
  * `secondaryText`
  * `tertiaryText`
  * `optionalText`
  * `size`
  * `presence`
  * `hideDetails`

* Removed `IFrame` properties:
  * `border`
* Removed `Stack` properties:
  * `border`
  * `borderLeft`
  * `borderRight`
  * `borderTop`
  * `borderBottom`
* Removed `Text` properties:
  * `borderLeft`
  * `borderRight`
  * `borderTop`
  * `borderBottom`

Bug fixes:

* Duplicate React rendering when loading a page.

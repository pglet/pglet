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
* New `page` properties:
  * `userAuthProvider`
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

Removed control properties:



Fixes:

* Duplicate React rendering when loading a page.

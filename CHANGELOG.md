# Change Log - Pglet Server

## [0.5.7](https://github.com/pglet/pglet/releases/tag/v0.5.7)

* `Stack` automatically scrolls to bottom if `autoscroll` property set to `true`.
* Set `page.UserAuthProvider` to a used authentication method (`github`, `google` or `azure`).
* `page.win_width` and `page.win_height` properties renamed to `page.winwidth` and `page.winheight`.
* When host is connected to a `page` its contents and properties are cleaned unless `update: true` is passed. No need to call `page.clean()` on the client anymore.
* Focusing input controls - allows setting focus on a control when added to a page or page loaded.
  * `Textbox.focused`
  * `Button.focused`
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

Fixes:

* Duplicate React rendering when loading a page.

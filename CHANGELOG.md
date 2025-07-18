# Changelog

## [v0.39.0](https://github.com/Songmu/deck/compare/v0.29.1...v0.39.0) - 2025-07-12

## [v0.39.0](https://github.com/k1LoW/deck/compare/v0.38.0...v0.39.0) - 2025-07-10
### New Features 🎉
- perf: add debounce function to avoid busy loop while watching by @Songmu in https://github.com/k1LoW/deck/pull/227
### Other Changes
- chore: fix lint errors and fix reviewdog/action-golangci-lint by @Songmu in https://github.com/k1LoW/deck/pull/229
- chore: set User-Agent to HTTP Request for crawling manner by @Songmu in https://github.com/k1LoW/deck/pull/228

## [v0.38.0](https://github.com/k1LoW/deck/compare/v0.37.1...v0.38.0) - 2025-07-09
### New Features 🎉
- feat(cmd): add markdown file support to deck new command by @Songmu in https://github.com/k1LoW/deck/pull/223
- perf: optimize watch process by @Songmu in https://github.com/k1LoW/deck/pull/226

## [v0.37.1](https://github.com/k1LoW/deck/compare/v0.37.0...v0.37.1) - 2025-07-08
### Fix bug 🐛
- fix: improve single-character list item rendering by @Songmu in https://github.com/k1LoW/deck/pull/220

## [v0.37.0](https://github.com/k1LoW/deck/compare/v0.36.0...v0.37.0) - 2025-07-07
### New Features 🎉
- feat: add support for skipping slides by @k1LoW in https://github.com/k1LoW/deck/pull/214
### Other Changes
- fix: rename JSON field names to use snake_case by @k1LoW in https://github.com/k1LoW/deck/pull/216

## [v0.36.0](https://github.com/k1LoW/deck/compare/v0.35.0...v0.36.0) - 2025-07-07
### New Features 🎉
- feat: Add ignore page functionality with page configuration comment by @Songmu in https://github.com/k1LoW/deck/pull/213
### Other Changes
- feat(md): add support for block quotes by @k1LoW in https://github.com/k1LoW/deck/pull/206
- chore(deps): bump the dependencies group with 2 updates by @dependabot in https://github.com/k1LoW/deck/pull/209

## [v0.35.0](https://github.com/k1LoW/deck/compare/v0.34.0...v0.35.0) - 2025-07-03
### New Features 🎉
- feat: Add {{output}} template variable support for code-block-to-image-command by @Songmu in https://github.com/k1LoW/deck/pull/204

## [v0.34.0](https://github.com/k1LoW/deck/compare/v0.33.0...v0.34.0) - 2025-07-02
### New Features 🎉
- make heading inside body bold; it's natural by @Songmu in https://github.com/k1LoW/deck/pull/196
### Other Changes
- doc: use camelCase for frontmatter fields by @Songmu in https://github.com/k1LoW/deck/pull/201

## [v0.33.0](https://github.com/k1LoW/deck/compare/v0.32.0...v0.33.0) - 2025-07-01
### New Features 🎉
- fix: refactor image handling for improved performance and flexibility by @k1LoW in https://github.com/k1LoW/deck/pull/191
- feat(apply): support setting title from frontmatter by @k1LoW in https://github.com/k1LoW/deck/pull/192
### Other Changes
- feat(cmd): add --presentation-id flag and deprecate positional arguments by @k1LoW in https://github.com/k1LoW/deck/pull/189

## [v0.32.0](https://github.com/k1LoW/deck/compare/v0.31.2...v0.32.0) - 2025-07-01
### New Features 🎉
- feat: add angle bracket autolinks support by @Songmu in https://github.com/k1LoW/deck/pull/185
- feat: support dynamic heading levels for title detection by @Songmu in https://github.com/k1LoW/deck/pull/188

## [v0.31.2](https://github.com/k1LoW/deck/compare/v0.31.1...v0.31.2) - 2025-07-01
### Fix bug 🐛
- fix: add CODEBLOCK_CONTENT env to cmd by @Songmu in https://github.com/k1LoW/deck/pull/181

## [v0.31.1](https://github.com/k1LoW/deck/compare/v0.31.0...v0.31.1) - 2025-07-01
### Fix bug 🐛
- fix: check sum comparison of images by @Songmu in https://github.com/k1LoW/deck/pull/179

## [v0.31.0](https://github.com/k1LoW/deck/compare/v0.30.0...v0.31.0) - 2025-06-30
### New Features 🎉
- feat: add YAML frontmatter support with new md.MD type by @Songmu in https://github.com/k1LoW/deck/pull/177
- feat: add presentationID support in frontmatter for simplified apply command by @Songmu in https://github.com/k1LoW/deck/pull/178

## [v0.30.0](https://github.com/k1LoW/deck/compare/v0.29.1...v0.30.0) - 2025-06-30
### New Features 🎉
- feat: create error.json if deck fails by @k1LoW in https://github.com/k1LoW/deck/pull/172
- feat: integrate tail-based logging and enhance error reporting by @k1LoW in https://github.com/k1LoW/deck/pull/174

## [v0.29.1](https://github.com/k1LoW/deck/compare/v0.29.0...v0.29.1) - 2025-06-30
### Other Changes
- refactor(deck): remove redundant DeletePageAfter call in ApplyPages function by @k1LoW in https://github.com/k1LoW/deck/pull/166
- refactor: prepend identifiers with type in object ID generation by @k1LoW in https://github.com/k1LoW/deck/pull/167
- feat(logger): add red color for failure messages by @k1LoW in https://github.com/k1LoW/deck/pull/169
- feat(apply): switch to JSON logging when `--verbose` and enhance action details structure by @k1LoW in https://github.com/k1LoW/deck/pull/170
- chore(deps): bump google.golang.org/api from 0.238.0 to 0.239.0 in the dependencies group by @dependabot in https://github.com/k1LoW/deck/pull/171

## [v0.29.0](https://github.com/k1LoW/deck/compare/v0.28.0...v0.29.0) - 2025-06-29
### New Features 🎉
- fix: prune old images via markdown by @k1LoW in https://github.com/k1LoW/deck/pull/165
### Fix bug 🐛
- fix: optimize markdown content comparisons by @k1LoW in https://github.com/k1LoW/deck/pull/162
- fix: respond when images are manually deleted in Google Slides by @k1LoW in https://github.com/k1LoW/deck/pull/164

## [v0.28.0](https://github.com/k1LoW/deck/compare/v0.27.1...v0.28.0) - 2025-06-29
### New Features 🎉
- feat: add global image caching mechanism by @k1LoW in https://github.com/k1LoW/deck/pull/160

## [v0.27.1](https://github.com/k1LoW/deck/compare/v0.27.0...v0.27.1) - 2025-06-29
### New Features 🎉
- fix: cache calculated checksums and pHash values for image comparison by @k1LoW in https://github.com/k1LoW/deck/pull/157
### Fix bug 🐛
- fix(md): correct image copying logic in ToSlides by @k1LoW in https://github.com/k1LoW/deck/pull/159

## [v0.27.0](https://github.com/k1LoW/deck/compare/v0.26.0...v0.27.0) - 2025-06-28
### New Features 🎉
- feat(dot): display a spinner when retrying a Google Slides API call by @k1LoW in https://github.com/k1LoW/deck/pull/151
- feat(logger): enhance dot handler for deleted/moved pages by @k1LoW in https://github.com/k1LoW/deck/pull/154
### Other Changes
- fix(deck): rename methods and add logging for page operations by @k1LoW in https://github.com/k1LoW/deck/pull/150
- fix: rename package by @k1LoW in https://github.com/k1LoW/deck/pull/152
- fix(deck): correct log levels for debug and warn messages by @k1LoW in https://github.com/k1LoW/deck/pull/153
- docs: correct `-l` option usage by @lacolaco in https://github.com/k1LoW/deck/pull/155
- feat(dot): enhance spinner handling by @k1LoW in https://github.com/k1LoW/deck/pull/156

## [v0.26.0](https://github.com/k1LoW/deck/compare/v0.25.0...v0.26.0) - 2025-06-26
### New Features 🎉
- feat: support bullet in shape by @k1LoW in https://github.com/k1LoW/deck/pull/147

## [v0.25.0](https://github.com/k1LoW/deck/compare/v0.24.2...v0.25.0) - 2025-06-24
### New Features 🎉
- feat(deck): add support for copying shapes between slides (create/delete) by @k1LoW in https://github.com/k1LoW/deck/pull/145
### Other Changes
- fix(md): rename `value` to `content` in code block by @k1LoW in https://github.com/k1LoW/deck/pull/144

## [v0.24.2](https://github.com/k1LoW/deck/compare/v0.24.1...v0.24.2) - 2025-06-23
### Fix bug 🐛
- fix(deck): correct indentation handling for nested paragraphs by @k1LoW in https://github.com/k1LoW/deck/pull/142

## [v0.24.1](https://github.com/k1LoW/deck/compare/v0.24.0...v0.24.1) - 2025-06-23
### New Features 🎉
- feat(deck): add functionality to copy images between (create/delete) slides when update by @k1LoW in https://github.com/k1LoW/deck/pull/139
### Fix bug 🐛
- fix(md): correct nested list parsing by @k1LoW in https://github.com/k1LoW/deck/pull/140
### Other Changes
- chore(deps): bump the dependencies group with 2 updates by @dependabot in https://github.com/k1LoW/deck/pull/137

## [v0.24.0](https://github.com/k1LoW/deck/compare/v0.23.0...v0.24.0) - 2025-06-22
### New Features 🎉
- feat(md): add support for code blocks in slide content by @k1LoW in https://github.com/k1LoW/deck/pull/129
- feat: add support for converting code blocks to images by @k1LoW in https://github.com/k1LoW/deck/pull/133
### Other Changes
- fix(md): introduce Parser struct for markdown parsing by @k1LoW in https://github.com/k1LoW/deck/pull/131
- Revert "fix(md): introduce Parser struct for markdown parsing" by @k1LoW in https://github.com/k1LoW/deck/pull/132
- fix(deck): simplify default*Layout logic by @k1LoW in https://github.com/k1LoW/deck/pull/134
- fix(md): use context with command by @k1LoW in https://github.com/k1LoW/deck/pull/135

## [v0.23.0](https://github.com/k1LoW/deck/compare/v0.22.1...v0.23.0) - 2025-06-21
### New Features 🎉
- feat(md): support image by @k1LoW in https://github.com/k1LoW/deck/pull/126
- feat(deck): support applying image ( `![img](path/to/image.png)` ) by @k1LoW in https://github.com/k1LoW/deck/pull/128

## [v0.22.1](https://github.com/k1LoW/deck/compare/v0.22.0...v0.22.1) - 2025-06-16
### Fix bug 🐛
- fix(deck): handle excess slides in ApplyPages by @k1LoW in https://github.com/k1LoW/deck/pull/124
### Other Changes
- chore(deps): bump google.golang.org/api from 0.236.0 to 0.237.0 in the dependencies group by @dependabot in https://github.com/k1LoW/deck/pull/123

## [v0.22.0](https://github.com/k1LoW/deck/compare/v0.21.1...v0.22.0) - 2025-06-15
### New Features 🎉
- fix: improve actions generation logic by @k1LoW in https://github.com/k1LoW/deck/pull/118
### Fix bug 🐛
- fix: improve action generation and similarity scoring by @k1LoW in https://github.com/k1LoW/deck/pull/121
### Other Changes
- test: add FuzzGenerateActions by @k1LoW in https://github.com/k1LoW/deck/pull/120
- test: add fuzzing workflow by @k1LoW in https://github.com/k1LoW/deck/pull/122

## [v0.21.1](https://github.com/k1LoW/deck/compare/v0.21.0...v0.21.1) - 2025-06-15
### Fix bug 🐛
- fix: resolve duplicate output in markdown lists separated by blank lines by @hanhan1978 in https://github.com/k1LoW/deck/pull/116
### Other Changes
- chore(deps): bump the dependencies group with 2 updates by @dependabot in https://github.com/k1LoW/deck/pull/113

## [v0.21.0](https://github.com/k1LoW/deck/compare/v0.20.0...v0.21.0) - 2025-06-12
### New Features 🎉
- fix: improve file change detection for vim compatibility by @hanhan1978 in https://github.com/k1LoW/deck/pull/114

## [v0.20.0](https://github.com/k1LoW/deck/compare/v0.19.1...v0.20.0) - 2025-06-08
### New Features 🎉
- fix: improve update logic for slides by @k1LoW in https://github.com/k1LoW/deck/pull/112
### Fix bug 🐛
- feat(deck): handle single-item bullet lists correctly by @k1LoW in https://github.com/k1LoW/deck/pull/111
### Other Changes
- chore(deps): bump google.golang.org/api from 0.234.0 to 0.235.0 in the dependencies group by @dependabot in https://github.com/k1LoW/deck/pull/108
- fix(deck): add integration test by @k1LoW in https://github.com/k1LoW/deck/pull/110

## [v0.19.1](https://github.com/k1LoW/deck/compare/v0.19.0...v0.19.1) - 2025-05-28
### Other Changes
- fix(md): handle allowed inline HTML elements by @k1LoW in https://github.com/k1LoW/deck/pull/107

## [v0.19.0](https://github.com/k1LoW/deck/compare/v0.18.0...v0.19.0) - 2025-05-28
### New Features 🎉
- feat(md): add support for class attributes in fragments by @k1LoW in https://github.com/k1LoW/deck/pull/102
- feat: support specifying styles for inline syntax in Markdown. by @k1LoW in https://github.com/k1LoW/deck/pull/104
- feat: add support for applying styles based on class names by @k1LoW in https://github.com/k1LoW/deck/pull/105
### Other Changes
- chore(deps): bump google.golang.org/api from 0.233.0 to 0.234.0 in the dependencies group by @dependabot in https://github.com/k1LoW/deck/pull/101

## [v0.18.0](https://github.com/k1LoW/deck/compare/v0.17.2...v0.18.0) - 2025-05-23
### New Features 🎉
- feat(oauth): implement PKCE for enhanced security in OAuth flow by @k1LoW in https://github.com/k1LoW/deck/pull/100
### Other Changes
- chore(deps): bump the dependencies group with 2 updates by @dependabot in https://github.com/k1LoW/deck/pull/99

## [v0.17.2](https://github.com/k1LoW/deck/compare/v0.17.1...v0.17.2) - 2025-05-18
### Fix bug 🐛
- fix: empty link item causes panic by @k1LoW in https://github.com/k1LoW/deck/pull/97
### Other Changes
- test: add fuzzing test by @k1LoW in https://github.com/k1LoW/deck/pull/95

## [v0.17.1](https://github.com/k1LoW/deck/compare/v0.17.0...v0.17.1) - 2025-05-18
### Fix bug 🐛
- fix: empty list item causes panic by @k1LoW in https://github.com/k1LoW/deck/pull/94

## [v0.17.0](https://github.com/k1LoW/deck/compare/v0.16.3...v0.17.0) - 2025-05-16
### New Features 🎉
- feat(deck): add support for custom code span styles in slides by @k1LoW in https://github.com/k1LoW/deck/pull/91

## [v0.16.3](https://github.com/k1LoW/deck/compare/v0.16.2...v0.16.3) - 2025-05-14
### Other Changes
- feat(md): add support for inline code fragments ( parse only ) by @k1LoW in https://github.com/k1LoW/deck/pull/88
- chore(deps): bump the dependencies group with 2 updates by @dependabot in https://github.com/k1LoW/deck/pull/86

## [v0.16.2](https://github.com/k1LoW/deck/compare/v0.16.1...v0.16.2) - 2025-05-05
### Fix bug 🐛
- fix(deck): count string as UTF-16 by @k1LoW in https://github.com/k1LoW/deck/pull/84

## [v0.16.1](https://github.com/k1LoW/deck/compare/v0.16.0...v0.16.1) - 2025-05-05
### Fix bug 🐛
- feat(deck): add emoji-aware string counting ( does not support ligature ) by @k1LoW in https://github.com/k1LoW/deck/pull/83
### Other Changes
- chore(deps): bump google.golang.org/api from 0.230.0 to 0.231.0 in the dependencies group by @dependabot in https://github.com/k1LoW/deck/pull/81

## [v0.16.0](https://github.com/k1LoW/deck/compare/v0.15.0...v0.16.0) - 2025-04-29
### New Features 🎉
- feat: add verbose flag and improve log message handling by @k1LoW in https://github.com/k1LoW/deck/pull/77
### Other Changes
- fix: use context.Context by @k1LoW in https://github.com/k1LoW/deck/pull/79

## [v0.15.0](https://github.com/k1LoW/deck/compare/v0.14.0...v0.15.0) - 2025-04-28
### New Features 🎉
- feat: support `--watch` for watching update file. by @k1LoW in https://github.com/k1LoW/deck/pull/76

## [v0.14.0](https://github.com/k1LoW/deck/compare/v0.13.5...v0.14.0) - 2025-04-28
### Other Changes
- refactor: Introduce a new `deck.Slide` type by @k1LoW in https://github.com/k1LoW/deck/pull/73

## [v0.13.5](https://github.com/k1LoW/deck/compare/v0.13.4...v0.13.5) - 2025-04-28
### Fix bug 🐛
- fix(deck): support additional title placeholder type by @k1LoW in https://github.com/k1LoW/deck/pull/72
### Other Changes
- chore(deps): bump the dependencies group with 2 updates by @dependabot in https://github.com/k1LoW/deck/pull/70

## [v0.13.4](https://github.com/k1LoW/deck/compare/v0.13.3...v0.13.4) - 2025-04-22
### Other Changes
- chore(deps): bump the dependencies group with 2 updates by @dependabot in https://github.com/k1LoW/deck/pull/68

## [v0.13.3](https://github.com/k1LoW/deck/compare/v0.13.2...v0.13.3) - 2025-04-17
### Other Changes
- chore(deps): bump golang.org/x/net from 0.37.0 to 0.38.0 by @dependabot in https://github.com/k1LoW/deck/pull/66

## [v0.13.2](https://github.com/k1LoW/deck/compare/v0.13.1...v0.13.2) - 2025-04-08
### Other Changes
- chore(deps): bump google.golang.org/api from 0.227.0 to 0.228.0 in the dependencies group by @dependabot in https://github.com/k1LoW/deck/pull/64
- chore(deps): bump golang.org/x/oauth2 from 0.28.0 to 0.29.0 in the dependencies group by @dependabot in https://github.com/k1LoW/deck/pull/65

## [v0.13.1](https://github.com/k1LoW/deck/compare/v0.13.0...v0.13.1) - 2025-03-30

## [v0.13.0](https://github.com/k1LoW/deck/compare/v0.12.0...v0.13.0) - 2025-03-30
### Other Changes
- fix: apply default layout only when creating a new one by @k1LoW in https://github.com/k1LoW/deck/pull/61

## [v0.12.0](https://github.com/k1LoW/deck/compare/v0.11.3...v0.12.0) - 2025-03-27
### New Features 🎉
- fix: support using refresh expired tokens by @k1LoW in https://github.com/k1LoW/deck/pull/58
### Other Changes
- chore(deps): bump google.golang.org/api from 0.226.0 to 0.227.0 in the dependencies group by @dependabot in https://github.com/k1LoW/deck/pull/57

## [v0.11.3](https://github.com/k1LoW/deck/compare/v0.11.2...v0.11.3) - 2025-03-22
### Fix bug 🐛
- fix(deck): clear placeholder for speaker notes in applyPage function by @k1LoW in https://github.com/k1LoW/deck/pull/55

## [v0.11.2](https://github.com/k1LoW/deck/compare/v0.11.1...v0.11.2) - 2025-03-20

## [v0.11.1](https://github.com/k1LoW/deck/compare/v0.11.0...v0.11.1) - 2025-03-20
### Fix bug 🐛
- fix: handle missing auth code in getTokenFromWeb function by @k1LoW in https://github.com/k1LoW/deck/pull/52

## [v0.11.0](https://github.com/k1LoW/deck/compare/v0.10.4...v0.11.0) - 2025-03-20
### New Features 🎉
- feat: support applying both bold and italic by @k1LoW in https://github.com/k1LoW/deck/pull/50

## [v0.10.4](https://github.com/k1LoW/deck/compare/v0.10.3...v0.10.4) - 2025-03-20
### Fix bug 🐛
- fix(deck): fix newline handling by @k1LoW in https://github.com/k1LoW/deck/pull/48

## [v0.10.3](https://github.com/k1LoW/deck/compare/v0.10.2...v0.10.3) - 2025-03-20
### Fix bug 🐛
- fix(deck): fix newline handling by @k1LoW in https://github.com/k1LoW/deck/pull/46

## [v0.10.2](https://github.com/k1LoW/deck/compare/v0.10.1...v0.10.2) - 2025-03-19
### Fix bug 🐛
- fix: remove unnecessary newline by @k1LoW in https://github.com/k1LoW/deck/pull/43
- fix(deck): clear text style before deleting paragraph bullets by @k1LoW in https://github.com/k1LoW/deck/pull/45

## [v0.10.1](https://github.com/k1LoW/deck/compare/v0.10.0...v0.10.1) - 2025-03-19
### Other Changes
- fix(md): trim whitespace and newline characters from HTML Block fragment values by @k1LoW in https://github.com/k1LoW/deck/pull/41

## [v0.10.0](https://github.com/k1LoW/deck/compare/v0.9.1...v0.10.0) - 2025-03-19
### New Features 🎉
- feat(deck): change sig and add start/end option by @k1LoW in https://github.com/k1LoW/deck/pull/38
### Other Changes
- fix: add `--page` option instead of `--start/end` by @k1LoW in https://github.com/k1LoW/deck/pull/40

## [v0.9.1](https://github.com/k1LoW/deck/compare/v0.9.0...v0.9.1) - 2025-03-18
### Fix bug 🐛
- fix(deck): reset bullet indices per body by @k1LoW in https://github.com/k1LoW/deck/pull/37

## [v0.9.0](https://github.com/k1LoW/deck/compare/v0.8.1...v0.9.0) - 2025-03-17
### New Features 🎉
- feat(md): add support for converting `<br>` tags to newlines by @k1LoW in https://github.com/k1LoW/deck/pull/33
### Other Changes
- chore(deps): bump the dependencies group with 2 updates by @dependabot in https://github.com/k1LoW/deck/pull/35

## [v0.8.1](https://github.com/k1LoW/deck/compare/v0.8.0...v0.8.1) - 2025-03-16
### Fix bug 🐛
- fix(md): remove empty bodies from parsed pages by @k1LoW in https://github.com/k1LoW/deck/pull/32

## [v0.8.0](https://github.com/k1LoW/deck/compare/v0.7.1...v0.8.0) - 2025-03-16
### New Features 🎉
- feat: add freeze functionality to skip page modifications by @k1LoW in https://github.com/k1LoW/deck/pull/29

## [v0.7.1](https://github.com/k1LoW/deck/compare/v0.7.0...v0.7.1) - 2025-03-16
### Other Changes
- fix(deck): increase retry limits for HTTP client by @k1LoW in https://github.com/k1LoW/deck/pull/27

## [v0.7.0](https://github.com/k1LoW/deck/compare/v0.6.1...v0.7.0) - 2025-03-16
### New Features 🎉
- feat(logging): integrate slog for enhanced logging by @k1LoW in https://github.com/k1LoW/deck/pull/26

## [v0.6.1](https://github.com/k1LoW/deck/compare/v0.6.0...v0.6.1) - 2025-03-15
### Fix bug 🐛
- fix(deck): correct text style update logic by @k1LoW in https://github.com/k1LoW/deck/pull/24

## [v0.6.0](https://github.com/k1LoW/deck/compare/v0.5.0...v0.6.0) - 2025-03-15
### New Features 🎉
- feat: add support for italic text by @k1LoW in https://github.com/k1LoW/deck/pull/22

## [v0.5.0](https://github.com/k1LoW/deck/compare/v0.4.1...v0.5.0) - 2025-03-15
### New Features 🎉
- feat(deck): add support for setting speaker notes by @k1LoW in https://github.com/k1LoW/deck/pull/20

## [v0.4.1](https://github.com/k1LoW/deck/compare/v0.4.0...v0.4.1) - 2025-03-15
### Other Changes
- fix: print to STDOUT by @k1LoW in https://github.com/k1LoW/deck/pull/18

## [v0.4.0](https://github.com/k1LoW/deck/compare/v0.3.0...v0.4.0) - 2025-03-15
### New Features 🎉
- feat: support link by @k1LoW in https://github.com/k1LoW/deck/pull/16

## [v0.3.0](https://github.com/k1LoW/deck/compare/v0.2.0...v0.3.0) - 2025-03-15
### New Features 🎉
- feat(apply): add `--title` flag to update presentation title by @k1LoW in https://github.com/k1LoW/deck/pull/13
- feat(cmd): add `new` command to create presentations by @k1LoW in https://github.com/k1LoW/deck/pull/14

## [v0.2.0](https://github.com/k1LoW/deck/compare/v0.1.1...v0.2.0) - 2025-03-14
### New Features 🎉
- feat: open browser automatically for OAuth authentication by @k1LoW in https://github.com/k1LoW/deck/pull/10
### Fix bug 🐛
- fix: Use context instead of channel by @k1LoW in https://github.com/k1LoW/deck/pull/8
### Other Changes
- fix(md): handle leading horizontal line in markdown by @k1LoW in https://github.com/k1LoW/deck/pull/11

## [v0.1.1](https://github.com/k1LoW/deck/compare/v0.1.0...v0.1.1) - 2025-03-14

## [v0.1.0](https://github.com/k1LoW/deck/compare/v0.0.1...v0.1.0) - 2025-03-14
### New Features 🎉
- feat: add `ls` command to list Google Slides presentations by @k1LoW in https://github.com/k1LoW/deck/pull/3
- feat(deck): add retryable HTTP client with exponential backoff by @k1LoW in https://github.com/k1LoW/deck/pull/5
### Other Changes
- chore: change directories by @k1LoW in https://github.com/k1LoW/deck/pull/6

## [v0.0.1](https://github.com/k1LoW/deck/commits/v0.0.1) - 2025-03-14
### Other Changes
- chore(deps): bump google.golang.org/api from 0.225.0 to 0.226.0 in the dependencies group by @dependabot in https://github.com/k1LoW/deck/pull/1

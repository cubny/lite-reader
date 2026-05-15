# mislav/go-readability

This is a fork of `github.com/go-shiori/go-readability` that adds performance optimizations, compatibility with [Readability.js 0.6.0](https://github.com/mozilla/readability/blob/main/CHANGELOG.md#060---2025-03-03), and some fixes for extracting article contents such as images.

## Changes

- Merge pull request #1 from mislav/dependabot/github_actions/all-ad97c67762 - readeck/go-readability@9b631b1
- Merge pull request #2 from mislav/dependabot/go_modules/all-136686d3d5 - readeck/go-readability@3e37e5f
- Merge pull request #3 from mislav/ci-tweaks - readeck/go-readability@82828cc
- Merge branch 'counting-optimizations' - readeck/go-readability@ef9026a
- Merge remote-tracking branch 'origin/master' - readeck/go-readability@95cc3ac
- Merge branch 'generate-test-timestamps' - readeck/go-readability@9ebfc8b
- Merge branch 'log-outerhtml' - readeck/go-readability@ad84ce8
- Merge branch 'preserve-div-ids' - readeck/go-readability@ee91de6
- Merge branch 'verbose-flag' - mislav/go-readability@f5ecec2
- Merge branch 'test-parser-truncate-diff' - readeck/go-readability@df47d01
- Merge branch 'figure-fix' - mislav/go-readability@a44e548
- Merge pull request #4 from mislav/linter-fixes - readeck/go-readability@31d0ef0
- Merge pull request #5 from mislav/parse-and-mutate - readeck/go-readability@daa20d1
- Merge pull request #8 from mislav/readability.js-0.6.0 - readeck/go-readability@772e5b1
- Merge pull request #11 from mislav/ci-rate-limit - readeck/go-readability@33469d9
- Merge pull request #10 from mislav/noscript-img - readeck/go-readability@460dbb5
- Merge pull request 'Rename module to codeberg.org/readeck/go-readability' (#2) from rename-module into main - readeck/go-readability@2b59a52
- Merge pull request 'go-readability CLI improvements' (#3) from cli-improvements into main - readeck/go-readability@ae7d57e
- Merge pull request 'Add configurable structured logger' (#4) from slog into main - readeck/go-readability@fb0fbc5
- Merge pull request 'Improve Article.TextContent rendering' (#5) from inner-text into main - readeck/go-readability@4c1efef
- Merge pull request 'Optimize searching for single element by tag name' (#6) from get-element into main - readeck/go-readability@0342e4f
- Merge pull request 'Avoid transforming \<br> chains inside \<pre> tags' (#7) from br-in-pre into main - readeck/go-readability@1c26ccc
- Merge pull request 'Simplify removing HTML comment nodes from the DOM' (#8) from cleanup-comment-nodes into main - readeck/go-readability@a8ff770
- Merge pull request 'Extend linting in CI to all modules in this project' (#9) from linting-fix into main - readeck/go-readability@d7abbe1
- Merge pull request 'Optimize cleanStyles' (#11) from clean-styles into main - readeck/go-readability@ea7e541
- Merge pull request 'Optimize common DOM operations' (#12) from optimize-traversal into main - readeck/go-readability@87ed180
- Merge remote-tracking branch 'origin/main' into content-debug - readeck/go-readability@d27f440
- Merge pull request 'Improve content-related logging' (#10) from content-debug into main - readeck/go-readability@c834769
- Merge pull request 'Reduce allocations when logging inside grabLoop' (#13) from logger-alloc-fix into main - readeck/go-readability@d48ccb1

### My pull requests to upstream

- [#100: Fix extracting images from figures on Medium](https://github.com/go-shiori/go-readability/pull/100)
- [#99: Optimize traversing the DOM when analyzing text content](https://github.com/go-shiori/go-readability/pull/99)
- [#98: Optimize logging HTML nodes by skipping costly compute when Debug is off](https://github.com/go-shiori/go-readability/pull/98)
- [#97: script/generate-test: add publishedTime, modifiedTime to expected metadata](https://github.com/go-shiori/go-readability/pull/97)
- [#96: Test_parser: make failures for text content mismatches more readable](https://github.com/go-shiori/go-readability/pull/96)
- [#95: Preserve IDs & class names of unwrapped DIVs](https://github.com/go-shiori/go-readability/pull/95)

## Benchmark

The benchmark measures the performance of parsing a very large HTML document (`test-pages/wikipedia-2/source.html`):

~~~
before: BenchmarkParser-8   	      24	  53734276 ns/op	73729978 B/op	  199153 allocs/op
after : BenchmarkParser-8   	      44	  25147170 ns/op	 5211741 B/op	   84818 allocs/op
~~~

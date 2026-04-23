# Changelog

## [0.5.1]

### Changed
- Removed escape system (defaults to fallthrough)
 

### Fixed
- RGBto256 fallback - catches more greys by using spread (max-min) instead of component differences between color channels


## [0.5.0]

### Added
- True color to 256 color palette fallback
- Changed syntax of light colors (lightred ==> lred)
- Added escape system `[<content>]`
- Dumb terminal detection

### Changed
- Hex validation now properly requires # prefix
- Fixed parse256ColorCode undefined variable bug
- Improved parseRGB length and int validation

### Fixed
- Unclosed bracket handling in templates
- ANSI 16 fallback logic


## [0.4.0] - 2026-04-06

### Added
- Inline padding in placeholders

---

## [0.3.0] - 2026-04-06

### Added
- Color support on windows

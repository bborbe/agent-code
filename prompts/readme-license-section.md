---
status: draft
created: "2026-08-09T16:35:00Z"
---

<summary>
- The project ships a license file but never mentions licensing in its readme
- A short license section is added at the end of the readme
- The wording matches the license file already in the repository
- No code changes and no new files
- Smallest possible change to satisfy the licensing convention
</summary>

<objective>
The readme ends with a License section that names the license and points at the existing license file.
</objective>

<context>
Read `CLAUDE.md` for project conventions if present (this repo has none; conventions come from sibling bborbe repos).

Files to read before making changes (read ALL first):
- `LICENSE` — already present at the repository root; a BSD-style license, copyright Benjamin Borbe
- `README.md` — currently ends with a bullet list of related paths and has no License section

The governing rule is `go-licensing/readme-license-section-required`. Check a sibling repository's readme (for example `github.com/bborbe/agent`) for the house wording before inventing your own.
</context>

<requirements>
1. Append a `## License` section as the last section of `README.md`.
2. Name the license as stated in the `LICENSE` file and link to that file.
3. Match the wording and heading level used by sibling bborbe repositories; do not invent a new format.
4. Do not alter the copyright years or holder — take them from `LICENSE` as-is.
5. Add a bullet under `## Unreleased` in `CHANGELOG.md` using a conventional prefix.
</requirements>

<constraints>
- Only change `README.md` and `CHANGELOG.md`
- Do NOT modify the `LICENSE` file
- Do NOT commit — dark-factory handles git
- Existing tests must still pass
</constraints>

<verification>
make precommit
</verification>

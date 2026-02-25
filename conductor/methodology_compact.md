# Conductor Methodology (Compact)

## English Overview (≈ 80 lines max)

1. **Kolmogorov Normalized Consolidation** – Treat every piece of project knowledge (specs, plans, code, tests) as a random variable in a shared probability space. By normalizing these variables we compute information entropy and enforce a minimal entropy increase for each new change, guaranteeing no hidden assumptions slip in.
2. **First‑Principles Decomposition** – Break any high‑level goal into atomic statements that are free of domain‑specific jargon. Each atom maps to a concrete markdown file in `conductor/` (e.g., `product.md`, `tech‑stack.md`), ensuring traceability from intent to implementation.
3. **Context‑Driven Development** – The AI agent never executes a task without loading **all** relevant markdown artifacts into its working context. This eliminates the need for external configuration files; the repo itself is the full spec.
4. **Agent‑Agnostic Protocols** – All commands (`setup`, `newTrack`, `implement`, …) are defined in TOML under `commands/conductor/`. Any AI coding agent that can read TOML can drive the workflow, making the system portable across Claude Code, OpenCode, Gemini CLI, Codex, etc.
5. **Iterative Entropy Check** – Before committing, the agent runs an entropy audit (`claude -p entropy-check`) which validates that added markdown/files do not increase the project's uncertainty beyond a defined threshold (default 0.01 bits per line).
6. **Bilingual Accessibility (optional)** – If needed, the Chinese version can be generated on‑demand via `claude -p translate` and cached as a token in a JSON file; the main markdown stays concise.
7. **Traceable Tokenisation** – Each paragraph in the English document can be assigned a unique token ID (e.g., `M001`). A JSON token file can map these IDs to translations, and agents can retrieve a token via `claude -p token-get M001`.
8. **Minimal LLM Context** – By keeping the methodology to ≤ 100 lines, we ensure a small LLM context while preserving full expressive power.

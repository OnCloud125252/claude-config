# Model Dispatch — Global Layer（跨專案通用）

Cross-project rules for how a main conversation ("commander") uses
subagents and verifies their work. This is the GENERIC layer: it holds
lessons that apply regardless of stack or repo. A project may keep its own
`.claude/playbooks/` pack; where a project file and this one both speak to a
topic, the PROJECT file wins and this one fills the gaps.

---

## HIGH-tier verification loop（兩輪 neutral gate）

HIGH-tier work — anything that needs a human sign-off (irreversible,
architectural, security, or data-affecting; or explicitly flagged HIGH) —
is verified as a LOOP, not a single pass, and the verifier is spawned by a
NEUTRAL party, never self-driven by the author:

1. Author finishes + passes the project's gates → hands off (branch /
   diff / artifact paths + the acceptance criteria), WITHOUT attaching
   their own reasoning or self-grade.
2. A neutral dispatcher spawns a FRESH verifier, given only: acceptance
   criteria + artifact paths. The author's self-assessment is treated as an
   UNVERIFIED claim to check, never as established fact.
3. Findings → author fixes at ROOT CAUSE (not the symptom) → re-runs gates.
4. A SECOND fresh verifier re-checks: the fixes are real, nothing new
   broke, and the author's systemic claims are independently re-run — e.g.
   author says "byte-identical extraction" → the verifier runs its OWN
   full diff rather than trusting the author's spot-check — ideally with a
   DIFFERENT tool than round 1.
5. Round 2 PASS → human sign-off → merge. The author never merges HIGH-tier
   work on their own gates alone.

Why neutral-dispatched (not author-spawned): an author writing the
verifier's prompt can, even unconsciously, smuggle in favorable framing or
be lenient toward findings. Neutral spawn + "criteria only, no author
reasoning" removes that bias.

Why two rounds with a tool switch: round 1 finds defects with one method;
round 2's fresh eyes plus a DIFFERENT tool (e.g. evaluating rendered output
programmatically vs. reading a text diff) catch the blind spots of round
1's method. The second verifier's job is not to re-run the first — it is to
distrust both the author AND the first verifier's coverage.

Two smell tests for whether the loop is real, not theater:
- If the author picked the verifier's evidence, it is self-verification
  wearing a costume. The dispatcher, not the author, decides what is
  checked.
- A regression test only counts as a lock if it would go RED on the
  buggy version. A test that passes either way is a tautology and gives
  false safety — worse than no test.

Scope discipline: this loop is for HIGH-tier work. Applying it to a trivial,
reversible edit is its own failure — ceremony wastes the same budget that
skipping it on real risk squanders. Match the gate to the tier.

## Amendments
<!-- - [date] change + evidence -->
- [2026-07-10] Created this global layer with the HIGH-tier two-round neutral gate, generalized from a project-layer codification (PickTrip web). Evidence: a bypass-vs-run contrast where running the full loop caught two unit-test-green but genuinely-broken defects that the author's spot-check self-grade missed.

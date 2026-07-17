# Sample HR Policy Questions

Use these questions to test the HR Policy Assistant from **Git Bash** on Windows:

```bash
go run ./cmd/ask -q "YOUR QUESTION HERE"
```

Or via the HTTP API (with `go run ./cmd/hrpolicy` running):

```bash
curl -X POST http://localhost:8081/api/query \
  -H "Content-Type: application/json" \
  -d '{"question": "YOUR QUESTION HERE"}'
```

---

## Leave Policy

- Can I take PL during notice period?
- How many privilege leave days do permanent employees get per year?
- Can I use casual leave during probation?
- What is the comp off validity period?
- How many days of paternity leave are available?
- What happens if I exhaust all leave and need more time off?
- How do I apply for sick leave for more than 2 consecutive days?

---

## Notice Period

- What is the notice period for an L3 engineer?
- What is the notice period during probation?
- Can the company put me on garden leave?
- Can I buy out my notice period?
- What happens if I don't serve the full notice period?
- Can I work from home during notice period?

---

## Exit & Separation

- What is the resignation process step by step?
- Can I withdraw my resignation after submitting it?
- When will I receive my full and final settlement?
- What assets must I return on my last day?
- What happens if an employee absconds?
- Is an exit interview mandatory?

---

## Work From Home / Hybrid

- How many WFH days are allowed per month?
- Who is not eligible for hybrid work?
- How do I apply for WFH?
- Can I work from home during notice period?
- What are the data security rules for WFH?
- Can I work from another country temporarily?

---

## Salary & Compensation

- When is salary credited each month?
- What are the components of CTC?
- How is annual increment decided?
- Is there a joining bonus and when is it paid?
- What reimbursements are employees eligible for?
- How long does full and final settlement take after exit?

---

## Onboarding & Entry

- What documents are required on joining day?
- How long is the probation period?
- What trainings are mandatory in the first 30 days?
- What is the buddy program?
- Can a former employee be rehired?
- What happens if I don't join on the agreed date?

---

## Cross-Policy (Advanced RAG Tests)

- I resigned with 30 days notice — can I take PL in my last week?
- What is the notice period for L2 vs L4, and can I WFH during that time?
- If I leave without serving notice, how does it affect my FnF?
- What leave can I take during probation and what is my notice period then?

---

## Quick Copy-Paste Commands

```bash
go run ./cmd/ask -q "What is the notice period for L3 engineer?"
go run ./cmd/ask -q "How many WFH days are allowed per month?"
go run ./cmd/ask -q "When is salary credited each month?"
go run ./cmd/ask -q "What is the full and final settlement timeline?"
go run ./cmd/ask -q "What documents are required on joining day?"
go run ./cmd/ask -q "Can I take comp off for working on a public holiday?"
go run ./cmd/ask -q "What happens if I abscond without serving notice?"
go run ./cmd/ask -q "Is garden leave paid and who decides it?"
```

---

## Policy Document Mapping

| Policy area        | Source file              |
| ------------------ | ------------------------ |
| Leave              | `documents/leave-policy.md`   |
| Notice period      | `documents/notice-policy.md`  |
| Exit & separation  | `documents/exit-policy.md`    |
| Work from home     | `documents/wfh-policy.md`     |
| Salary             | `documents/salary-policy.md`  |
| Onboarding & entry | `documents/entry-policy.md`   |

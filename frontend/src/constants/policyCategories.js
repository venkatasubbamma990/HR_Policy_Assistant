export const POLICY_CATEGORIES = [
  {
    id: 'leave',
    label: 'Leave Policy',
    shortLabel: 'Leave',
    icon: 'calendar',
    questions: [
      'Can I take PL during notice period?',
      'How many privilege leave days do permanent employees get?',
      'How many days of paternity leave are available?',
    ],
  },
  {
    id: 'notice',
    label: 'Notice Period',
    shortLabel: 'Notice',
    icon: 'clock',
    questions: [
      'What is the notice period for an L3 engineer?',
      'Can I buy out my notice period?',
      'Can I work from home during notice period?',
    ],
  },
  {
    id: 'wfh',
    label: 'Work From Home',
    shortLabel: 'WFH',
    icon: 'home',
    questions: [
      'How many WFH days are allowed per month?',
      'Who is not eligible for hybrid work?',
      'What are the data security rules for WFH?',
    ],
  },
  {
    id: 'salary',
    label: 'Salary & Compensation',
    shortLabel: 'Salary',
    icon: 'wallet',
    questions: [
      'When is salary credited each month?',
      'What are the components of CTC?',
      'How long does full and final settlement take?',
    ],
  },
  {
    id: 'exit',
    label: 'Exit & Separation',
    shortLabel: 'Exit',
    icon: 'logout',
    questions: [
      'What is the resignation process step by step?',
      'When will I receive my full and final settlement?',
      'What happens if an employee absconds?',
    ],
  },
  {
    id: 'entry',
    label: 'Onboarding',
    shortLabel: 'Onboarding',
    icon: 'user',
    questions: [
      'What documents are required on joining day?',
      'How long is the probation period?',
      'What trainings are mandatory in the first 30 days?',
    ],
  },
]

export const DEFAULT_CATEGORY = POLICY_CATEGORIES[0]

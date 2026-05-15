export type LegalDocumentKind = 'terms' | 'privacy'

export interface LegalSection {
  title: string
  paragraphs?: string[]
  bullets?: string[]
  subsections?: LegalSection[]
}

export interface LegalDocument {
  kind: LegalDocumentKind
  title: string
  eyebrow: string
  summary: string
  readingNote: string
  lastUpdated: string
  version: string
  sections: LegalSection[]
}

const providerPolicyLinks = [
  'OpenAI: https://openai.com/policies/privacy-policy',
  'Anthropic: https://www.anthropic.com/legal/privacy',
  'Google: https://policies.google.com/privacy',
  'DeepSeek: https://cdn.deepseek.com/policies/en-US/deepseek-privacy-policy.html',
  'Moonshot AI: https://platform.moonshot.ai/privacy',
  'Zhipu AI: https://z.ai/privacy'
]

export const termsDocument: LegalDocument = {
  kind: 'terms',
  title: 'User Agreement',
  eyebrow: 'Service rules and responsibilities',
  summary:
    'Please read this agreement carefully before using cCoder.me, especially account security, service limits, billing rules, user content, third-party AI model terms, liability boundaries, and agreement changes.',
  readingNote:
    'Pay special attention to bold headings, clause lists, and the body content configured by the administrator. If any summary differs from the body, the body prevails.',
  lastUpdated: 'May 15, 2026',
  version: '2026-05-15',
  sections: [
    {
      title: '1. Overview of cCoder.me Services',
      paragraphs: [
        'This User Agreement (this "Agreement") is a legal agreement between you and cCoder.me, its owner, and authorized operators ("cCoder.me," "we," "us," or "our") regarding your access to and use of https://ccoder.me, the web console, APIs, documentation, and related services (collectively, the "Services").',
        'cCoder.me operates an AI API gateway and developer platform. The Services may allow users to create API keys, route requests to third-party artificial intelligence model providers, manage subscriptions or usage balance, view usage records, and access documentation and support resources.',
        'The Services rely on third-party model providers, infrastructure providers, payment processors, authentication providers, and other external services. We may add, remove, suspend, or change models, routes, prices, quotas, or features at any time.'
      ]
    },
    {
      title: '2. Eligibility and Scope of Users',
      bullets: [
        'You must be at least 18 years old, or the age of legal majority in your jurisdiction, to use the Services.',
        'If you use the Services on behalf of a company, organization, or other legal entity, you represent that you have authority to bind that entity to this Agreement.',
        'You may not use the Services if applicable law prohibits you from doing so, if you are located in a restricted jurisdiction, or if your account has been suspended or terminated by us.',
        'You are responsible for complying with all laws and rules that apply to your location, your users, your applications, your prompts, your outputs, and your use of third-party AI models.'
      ]
    },
    {
      title: '3. Account and Registration',
      bullets: [
        'Most Services require an account. You agree to provide accurate, complete, and current registration information.',
        'You are responsible for maintaining the confidentiality of your password, API keys, OAuth sessions, and other credentials.',
        'You are responsible for all activity under your account, including charges caused by leaked API keys or credentials.',
        'You must promptly notify us through the support contact published on the platform if you suspect unauthorized access or credential compromise.',
        'We may reject, suspend, or terminate accounts that use false information, evade limits, abuse promotions, threaten platform integrity, or violate this Agreement.'
      ]
    },
    {
      title: '4. Payment, Credits, Subscriptions, and Refunds',
      paragraphs: [
        'The Services may require prepaid balance, credits, subscription purchases, quota packages, or other payments. Prices, included quota, validity periods, renewal rules, and supported payment methods are shown in the console or purchase flow at the time of purchase.',
        'Unless mandatory law requires otherwise, purchased credits, subscriptions, quota packages, and other digital services are non-refundable once delivered or consumed. We may provide refunds when we discontinue a paid service, when a material service failure caused by us prevents you from using purchased unused credits, or when applicable law requires a refund.',
        'Taxes, bank fees, currency conversion fees, payment processor fees, and third-party charges may apply. You authorize us and our payment processors to charge your selected payment method for amounts you approve.'
      ]
    },
    {
      title: '5. API Use, Quotas, and Service Limits',
      bullets: [
        'You must use the API only through documented endpoints and supported authentication methods.',
        'We may enforce request limits, concurrency limits, rate limits, spend limits, model availability limits, abuse filters, and other operational controls.',
        'Usage records and billing calculations are based on our metering systems and third-party provider usage data. We may correct billing errors when detected.',
        'You may not bypass technical limits, share credentials for resale, overload the Services, scrape private endpoints, or interfere with platform stability.',
        'We do not guarantee that any model, route, context window, output format, latency, price, or feature will remain available.'
      ]
    },
    {
      title: '6. User Content and AI Outputs',
      paragraphs: [
        'You may submit prompts, code, files, images, text, metadata, and other materials to the Services ("Inputs") and receive model-generated responses ("Outputs"). Inputs and Outputs are collectively "User Content."',
        'You retain any rights you have in your Inputs. Rights in Outputs are governed by applicable law and the terms of the model provider that generated them. We do not guarantee that Outputs are unique, accurate, lawful, non-infringing, secure, or suitable for any purpose.',
        'You are solely responsible for evaluating Outputs before relying on them. You should not use Outputs as a substitute for professional advice in legal, medical, financial, safety-critical, or other high-risk contexts.'
      ]
    },
    {
      title: '7. Third-Party AI Model Terms and Training',
      paragraphs: [
        'Your use of a model may be subject to the model provider\'s own terms, policies, usage rules, safety rules, data handling practices, logging practices, and training choices. You are responsible for reviewing and complying with those terms.',
        'To the extent technically and commercially available, we may configure provider accounts to reduce or disable model training on customer data. We do not control all third-party providers and cannot guarantee how every provider handles Inputs, Outputs, logs, abuse review, or training.'
      ],
      bullets: providerPolicyLinks
    },
    {
      title: '8. Private Logging and Debugging',
      paragraphs: [
        'The Services may offer request logging, usage records, debugging views, private input/output logging, or similar features. If you enable such features, you authorize us to store, reproduce, process, and display related User Content solely to provide those features, operate the Services, investigate abuse, debug issues, comply with law, and protect the Services.',
        'You should not submit sensitive personal data, regulated health data, payment card data, secrets, private keys, or confidential third-party data unless you have the rights, consents, and security controls required for that data.'
      ]
    },
    {
      title: '9. Input Representations and Warranties',
      bullets: [
        'You represent that you own or have all rights, licenses, consents, permissions, and legal bases required to submit Inputs and use Outputs.',
        'Your Inputs and use of Outputs must not violate law, infringe intellectual property rights, violate privacy or publicity rights, breach confidentiality duties, or cause us or a model provider to violate applicable rules.',
        'You are responsible for your users and applications that access the Services through your account or API keys.'
      ]
    },
    {
      title: '10. Prohibited Conduct',
      bullets: [
        'Use the Services for unlawful, harmful, deceptive, abusive, harassing, exploitative, or infringing activity.',
        'Use the Services for political campaigning, election interference, voter persuasion, public-opinion manipulation, political propaganda, lobbying automation, targeted political messaging, or generation or distribution of political content intended to influence civic processes.',
        'Use the Services to create, optimize, distribute, or conceal content that incites violence, terrorism, extremism, self-harm, sexual exploitation, human trafficking, fraud, money laundering, sanctions evasion, illegal gambling, controlled-substance transactions, or other illegal activity.',
        'Generate or distribute malware, credential theft tools, exploit instructions, spam, phishing content, or content that meaningfully facilitates cyber abuse.',
        'Use the Services to evade law enforcement, regulatory oversight, platform safety systems, sanctions, export controls, identity verification, content moderation, usage limits, billing controls, or abuse detection.',
        'Submit, process, or attempt to extract highly sensitive personal data, credentials, trade secrets, private keys, classified information, or regulated data without all required authority, notices, consents, safeguards, and legal bases.',
        'Generate impersonation, deepfake, deceptive, defamatory, discriminatory, hateful, violent, sexual, exploitative, or otherwise abusive content, or content that violates the rights, safety, privacy, or dignity of any person or group.',
        'Attempt to reverse engineer, bypass, disable, overload, scrape, probe, or interfere with the Services, security controls, billing systems, model routing, or provider integrations.',
        'Create multiple accounts or false identities to evade limits, sanctions, billing, suspension, review, or abuse controls.',
        'Resell, sublicense, broker, or redistribute API access unless we have authorized you in writing.',
        'Use the Services to violate third-party model provider terms, safety policies, export controls, sanctions, or applicable law.',
        'Submit personal data or confidential information without all required notices, consents, authority, and safeguards.',
        'Assist, encourage, or permit anyone else to do any prohibited act.'
      ]
    },
    {
      title: '11. Suspension, Termination, and Changes to Services',
      paragraphs: [
        'You may stop using the Services at any time. You remain responsible for charges incurred before termination.',
        'We may suspend, restrict, or terminate your access at any time if we believe you violated this Agreement, created security or legal risk, failed to pay, abused the Services, or if continued service is commercially or technically impractical.',
        'We may modify, discontinue, or deprecate any feature, model, endpoint, pricing plan, payment method, or documentation. We will try to provide notice of material changes when practical, but urgent security, legal, provider, or operational changes may take effect immediately.'
      ]
    },
    {
      title: '12. Privacy Policy',
      paragraphs: [
        'Our Privacy Policy explains how we collect, use, disclose, retain, and protect personal data. The Privacy Policy is incorporated into this Agreement by reference. By using the Services, you acknowledge our personal data practices described in the Privacy Policy.'
      ]
    },
    {
      title: '13. Ownership and Proprietary Rights',
      paragraphs: [
        'The Services, software, interfaces, dashboards, designs, documentation, logos, trademarks, model routing systems, billing systems, and other platform materials are owned by us or our licensors and are protected by intellectual property and other laws.',
        'Except for the limited right to use the Services under this Agreement, no rights are transferred to you. You may not copy, modify, distribute, sell, lease, or create derivative works of the Services unless we expressly authorize it.'
      ]
    },
    {
      title: '14. Feedback',
      paragraphs: [
        'If you provide comments, ideas, bug reports, feature requests, or other feedback, you grant us a perpetual, irrevocable, worldwide, royalty-free license to use that feedback for any purpose, including improving, developing, marketing, and operating the Services, without compensation to you.'
      ]
    },
    {
      title: '15. Confidentiality',
      paragraphs: [
        'If either party receives non-public information from the other party that is identified as confidential or should reasonably be understood as confidential, the receiving party must use reasonable care to protect it and use it only for purposes related to the Services.',
        'Confidential information does not include information that is public without breach, already known without confidentiality duty, independently developed, lawfully received from a third party, or feedback. We may disclose confidential information when required by law, court order, regulator request, or to protect rights, safety, security, or service integrity.'
      ]
    },
    {
      title: '16. Indemnification',
      paragraphs: [
        'You agree to defend, indemnify, and hold harmless cCoder.me, its operators, affiliates, contractors, service providers, officers, employees, and agents from claims, damages, liabilities, losses, costs, and expenses, including reasonable attorneys\' fees, arising from your use of the Services, User Content, breach of this Agreement, violation of law, infringement of rights, or dispute with a third party.'
      ]
    },
    {
      title: '17. Disclaimer; No Warranties',
      paragraphs: [
        'THE SERVICES, MODELS, OUTPUTS, DOCUMENTATION, AND ALL MATERIALS ARE PROVIDED ON AN "AS IS" AND "AS AVAILABLE" BASIS. TO THE MAXIMUM EXTENT PERMITTED BY LAW, WE DISCLAIM ALL WARRANTIES, EXPRESS OR IMPLIED, INCLUDING WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, TITLE, QUIET ENJOYMENT, NON-INFRINGEMENT, ACCURACY, AVAILABILITY, SECURITY, AND ERROR-FREE OPERATION.',
        'WE DO NOT WARRANT THAT THE SERVICES WILL BE UNINTERRUPTED, SECURE, TIMELY, ACCURATE, OR FREE OF ERRORS, VIRUSES, HARMFUL COMPONENTS, PROVIDER FAILURES, MODEL CHANGES, OR DATA LOSS.',
        'WE DO NOT CONTROL THIRD-PARTY MODEL PROVIDERS, PAYMENT PROCESSORS, INFRASTRUCTURE PROVIDERS, NETWORK PROVIDERS, OR USER APPLICATIONS. WE ARE NOT RESPONSIBLE FOR THEIR ACTS, OMISSIONS, TERMS, PRICES, RATE LIMITS, SAFETY FILTERS, TRAINING PRACTICES, DATA HANDLING, AVAILABILITY, OUTPUT QUALITY, OR LEGAL COMPLIANCE.',
        'YOU ARE SOLELY RESPONSIBLE FOR PROMPTS, INPUTS, OUTPUTS, INTEGRATIONS, APPLICATIONS, USERS, DEPLOYMENTS, DECISIONS, AND ANY CONSEQUENCES OF USING OR RELYING ON AI-GENERATED CONTENT. OUTPUTS MAY BE FALSE, INCOMPLETE, BIASED, UNSAFE, OFFENSIVE, OUTDATED, OR INFRINGING.'
      ]
    },
    {
      title: '18. Limitation of Liability',
      paragraphs: [
        'TO THE MAXIMUM EXTENT PERMITTED BY LAW, WE WILL NOT BE LIABLE FOR INDIRECT, INCIDENTAL, SPECIAL, CONSEQUENTIAL, EXEMPLARY, OR PUNITIVE DAMAGES, INCLUDING LOST PROFITS, LOST REVENUE, LOST GOODWILL, LOST DATA, BUSINESS INTERRUPTION, SECURITY INCIDENTS CAUSED BY YOUR CREDENTIALS, OR DAMAGES ARISING FROM OUTPUTS OR THIRD-PARTY PROVIDERS.',
        'TO THE MAXIMUM EXTENT PERMITTED BY LAW, OUR TOTAL LIABILITY FOR ALL CLAIMS ARISING OUT OF OR RELATING TO THE SERVICES OR THIS AGREEMENT WILL NOT EXCEED THE GREATER OF (A) THE AMOUNT YOU PAID TO US FOR THE SERVICES IN THE 12 MONTHS BEFORE THE EVENT GIVING RISE TO LIABILITY, OR (B) US$100.',
        'THESE LIMITATIONS APPLY EVEN IF A REMEDY FAILS OF ITS ESSENTIAL PURPOSE, EVEN IF WE WERE ADVISED OF THE POSSIBILITY OF DAMAGES, AND REGARDLESS OF THE THEORY OF LIABILITY, INCLUDING CONTRACT, TORT, NEGLIGENCE, STRICT LIABILITY, STATUTE, OR OTHERWISE.'
      ]
    },
    {
      title: '19. User Acknowledgement for Registration and API Key Use',
      paragraphs: [
        'Before creating an account or API key, you must affirmatively confirm that you have read, understood, and accepted this Agreement and the Privacy Policy. You must also confirm that you will not use the Services for political activity, illegal activity, abuse, evasion, infringement, or any other prohibited conduct.',
        'If you do not provide this confirmation, we may reject registration, API key creation, or API requests. If you breach this confirmation, we may suspend or terminate your account, revoke API keys, preserve relevant records, cooperate with lawful requests, and pursue any available remedy.'
      ]
    },
    {
      title: '20. Governing Law and Dispute Resolution',
      paragraphs: [
        'Unless a mandatory law in your jurisdiction requires otherwise, this Agreement is governed by the laws of Hong Kong, excluding conflict-of-law rules. The parties will first attempt to resolve disputes in good faith through written notice and negotiation.',
        'If a dispute cannot be resolved within 30 days after notice, either party may submit the dispute to the competent courts of Hong Kong, unless the parties separately agree to arbitration or another forum in writing. Either party may seek injunctive or equitable relief for intellectual property, security, confidentiality, or unauthorized access matters in any court with jurisdiction.'
      ]
    },
    {
      title: '21. General Provisions and Contact',
      paragraphs: [
        'This Agreement, the Privacy Policy, and policies incorporated by reference constitute the entire agreement between you and us regarding the Services. If any provision is unenforceable, the remaining provisions remain in effect. Our failure to enforce a provision is not a waiver.',
        'You may not assign this Agreement without our prior written consent. We may assign this Agreement in connection with a merger, acquisition, reorganization, sale of assets, change of control, or by operation of law.',
        'Questions about this Agreement may be sent to the support contact published on cCoder.me or in your account console.'
      ]
    }
  ]
}

export const privacyDocument: LegalDocument = {
  kind: 'privacy',
  title: 'Privacy Policy',
  eyebrow: 'Data protection and privacy statement',
  summary:
    'We value user privacy and data security. This policy explains how cCoder.me collects, uses, stores, protects, shares, and otherwise processes information.',
  readingNote:
    'Pay special attention to bold headings, clause lists, and the body content configured by the administrator. If any summary differs from the body, the body prevails.',
  lastUpdated: 'May 15, 2026',
  version: '2026-05-15',
  sections: [
    {
      title: '1. Personal Data We Collect',
      paragraphs: [
        'This Privacy Policy applies to personal data collected when you visit cCoder.me, create an account, use the console, call the API, interact with documentation or support, make purchases, or otherwise use the Services.',
        'Personal data means information that identifies, relates to, describes, or can reasonably be linked to an individual. We may collect personal data directly from you, automatically from your use of the Services, and from third parties such as authentication providers, payment processors, model providers, anti-abuse vendors, and infrastructure providers.'
      ],
      subsections: [
        {
          title: '1.1 Personal data you provide',
          bullets: [
            'Account information, such as email address, username, password hash, authentication method, profile settings, and support contact details.',
            'Payment and billing information, such as order records, recharge records, invoices, payment method metadata, transaction identifiers, and payment processor responses.',
            'User Content, including prompts, files, code, images, metadata, API request bodies, and model outputs if you submit them or enable logging features.',
            'Communications, such as support messages, feedback, dispute notices, and other correspondence.'
          ]
        },
        {
          title: '1.2 Personal data collected automatically',
          bullets: [
            'Device and network data, such as IP address, browser type, operating system, user agent, language settings, time zone, device identifiers, referrer URLs, and approximate location inferred from IP address.',
            'Usage and log data, such as login times, API endpoints called, model names, token or usage counts, latency, errors, status codes, request IDs, balance changes, subscription activity, and security events.',
            'Cookie and local storage data used to keep you signed in, remember settings, secure sessions, prevent abuse, and understand product usage.'
          ]
        },
        {
          title: '1.3 Third-party model provider processing',
          paragraphs: [
            'When you route requests to third-party AI model providers, those providers may process Inputs, Outputs, metadata, logs, and abuse signals under their own terms and privacy policies. We do not fully control their processing, retention, or training practices.'
          ],
          bullets: providerPolicyLinks
        }
      ]
    },
    {
      title: '2. How We Use Personal Data',
      bullets: [
        'Provide, operate, authenticate, secure, maintain, troubleshoot, and improve the Services.',
        'Create accounts, verify identity, manage sessions, issue API keys, route model requests, calculate usage, process billing, and provide customer support.',
        'Detect, prevent, investigate, and respond to fraud, abuse, credential compromise, spam, malware, excessive load, policy violations, and security incidents.',
        'Send service notices, account notices, security alerts, payment notices, policy updates, support responses, and administrative messages.',
        'Analyze service performance, model routing, feature adoption, errors, and aggregate usage trends.',
        'Comply with legal obligations, enforce agreements, resolve disputes, protect rights and safety, and respond to lawful requests.'
      ]
    },
    {
      title: '3. Cookies and Similar Technologies',
      paragraphs: [
        'We use cookies, local storage, session storage, pixels, logs, and similar technologies to operate and secure the Services, remember preferences, maintain sessions, measure product performance, and analyze usage.',
        'You can configure your browser to block or delete cookies. If you disable required cookies or storage, some parts of the Services may not function properly, including login, security, billing, or console features.'
      ],
      subsections: [
        {
          title: '3.1 Required technologies',
          paragraphs: [
            'Required cookies and storage are necessary for authentication, security, fraud prevention, load balancing, language settings, and core product functionality. These cannot be disabled through the Services.'
          ]
        },
        {
          title: '3.2 Functional and analytics technologies',
          paragraphs: [
            'Functional and analytics technologies help us understand usage, improve product quality, debug errors, measure feature adoption, and provide a better user experience. Where required by law, we will request consent or provide opt-out choices.'
          ]
        }
      ]
    },
    {
      title: '4. How We Share and Disclose Personal Data',
      bullets: [
        'Service providers and processors: infrastructure hosting, databases, email delivery, security, analytics, payment processing, customer support, and other vendors that process data for us.',
        'Third-party AI model providers: providers that receive API requests, Inputs, Outputs, metadata, and abuse signals when you use their models through the Services.',
        'Payment processors and financial partners: providers that process payments, refunds, fraud checks, tax records, and transaction disputes.',
        'Authentication and integration providers: OAuth, single sign-on, notification, or platform integration providers you choose to use.',
        'Affiliates, contractors, and professional advisers: parties that support operations, legal compliance, accounting, auditing, security, and business administration.',
        'Legal and safety recipients: regulators, law enforcement, courts, public authorities, or other third parties when required by law or when we believe disclosure is reasonably necessary to protect rights, safety, security, or service integrity.',
        'Business transaction recipients: counterparties, advisers, and successors in connection with a merger, financing, acquisition, reorganization, bankruptcy, sale of assets, or change of control.',
        'With your consent or direction: any other sharing you authorize or request.'
      ]
    },
    {
      title: '5. Aggregated and De-Identified Data',
      paragraphs: [
        'We may process personal data into aggregated, anonymized, or de-identified information that does not reasonably identify you. We may use and share such information for analytics, reporting, product improvement, security research, pricing, capacity planning, and business purposes.',
        'If we maintain de-identified data, we will take reasonable measures designed to prevent re-identification except as permitted by law, such as testing whether de-identification is effective.'
      ]
    },
    {
      title: '6. Data Security',
      paragraphs: [
        'We use administrative, technical, and organizational measures designed to protect personal data from unauthorized access, loss, misuse, alteration, and disclosure. Measures may include encryption in transit, access controls, audit logs, credential hashing, network protections, backups, and monitoring.',
        'No internet service is completely secure. You are responsible for protecting your account password, API keys, OAuth sessions, devices, and applications. If you believe your account or API key has been compromised, contact support and rotate credentials immediately.'
      ]
    },
    {
      title: '7. Data Retention',
      paragraphs: [
        'We retain personal data for as long as reasonably necessary to provide the Services, maintain accounts, process payments, comply with legal obligations, resolve disputes, enforce agreements, prevent abuse, maintain security, and support legitimate business needs.',
        'Retention periods vary by data type. Account and billing records may be retained for legal, tax, accounting, and audit purposes. Security logs may be retained to investigate abuse and protect the platform. User Content logging depends on product settings, debugging needs, and applicable provider practices.',
        'When personal data is no longer needed, we will delete, anonymize, or otherwise dispose of it according to applicable law and operational requirements.'
      ]
    },
    {
      title: '8. Your Rights and Choices',
      bullets: [
        'Access: you may request a copy of personal data we process about you.',
        'Correction: you may request correction of inaccurate or incomplete personal data.',
        'Deletion: you may request deletion of personal data, subject to legal, security, billing, fraud-prevention, and operational exceptions.',
        'Portability: where required by law, you may request export of certain personal data in a portable format.',
        'Restriction or objection: where required by law, you may object to or request restriction of certain processing.',
        'Consent withdrawal: where processing is based on consent, you may withdraw consent at any time, without affecting processing that occurred before withdrawal.',
        'Cookie choices: you may control browser cookies and storage, subject to required technologies needed for the Services.'
      ],
      paragraphs: [
        'To exercise rights, use account settings where available or contact support through the support contact published on cCoder.me or in your account console. We may need to verify your identity before acting on a request.'
      ]
    },
    {
      title: '9. International Data Transfers',
      paragraphs: [
        'We and our service providers may process personal data in countries other than where you live. These countries may have data protection laws different from those in your jurisdiction.',
        'When required, we rely on appropriate transfer mechanisms such as adequacy decisions, standard contractual clauses, data processing agreements, consent, contractual necessity, or other mechanisms recognized by applicable law.'
      ]
    },
    {
      title: '10. Children and Age Restrictions',
      paragraphs: [
        'The Services are not directed to children and are intended only for users who are at least 18 years old or the age of legal majority in their jurisdiction. We do not knowingly collect personal data from children. If you believe a child provided personal data to us, contact support so we can review and take appropriate action.'
      ]
    },
    {
      title: '11. Supplemental Disclosures for GDPR, UK GDPR, and U.S. State Privacy Laws',
      paragraphs: [
        'Where the GDPR, UK GDPR, or similar laws apply, our legal bases may include performance of a contract, legitimate interests, consent, compliance with legal obligations, and protection of rights and safety. Our legitimate interests include operating and securing the Services, preventing abuse, improving products, communicating with users, and enforcing agreements.',
        'Where U.S. state privacy laws apply, categories of personal data we may collect include identifiers, commercial information, internet or network activity, approximate geolocation, account credentials, communications, and inferences related to product usage. We use and disclose these categories for the purposes described in this Policy.',
        'We do not sell personal data for money. We also do not knowingly sell or share personal data of children. Some analytics or advertising technologies, if enabled in the future, may be considered "sharing" or "targeted advertising" under certain laws; where required, we will provide applicable notices and opt-out choices.',
        'You may have the right not to be discriminated against for exercising privacy rights. Authorized agents may submit requests where permitted by law, but we may require proof of authorization and identity verification.'
      ]
    },
    {
      title: '12. Changes and Contact',
      paragraphs: [
        'We may update this Privacy Policy from time to time. If changes are material, we will provide notice by posting the updated policy, emailing account holders when practical, or using another reasonable method. Your continued use of the Services after the effective date means you acknowledge the updated policy.',
        'Questions about this Privacy Policy or privacy requests may be sent to the support contact published on cCoder.me or in your account console.'
      ]
    }
  ]
}

export const legalDocuments: Record<LegalDocumentKind, LegalDocument> = {
  terms: termsDocument,
  privacy: privacyDocument
}

export function getLegalDocument(kind: LegalDocumentKind): LegalDocument {
  return legalDocuments[kind]
}

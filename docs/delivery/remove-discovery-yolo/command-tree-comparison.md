# Command Tree Comparison

- generated_at: `2026-07-07T06:11:35Z`
- scope: visible `--help` command tree only; hidden compatibility commands are excluded by design.
- open_delivery_commit: `1b823514`
- wukong_local_head: `f48ae79d` with local compatibility bridge changes; `origin/develop` (`2a7c70ff`) currently fails to compile against the remove-discovery core because it still references removed `pkg/cmdutil` overlay APIs.

| tree | visible commands | help errors |
| --- | ---: | ---: |
| open-delivery | 779 | 0 |
| wukong-local | 973 | 0 |

## Top-Level Difference

Wukong-only top-level commands:
- `agoal`
- `aiapp`
- `aidesign`
- `blackboard`
- `conference`
- `credit`
- `docparse`
- `finance`
- `law`
- `yida`

Open-only top-level commands:
- `doctor`
- `upgrade`

## Full Visible Path Difference

- wukong_only_visible_paths: `200`
- open_only_visible_paths: `6`

See:
- `command-tree-wukong-local.md` / `.json`
- `command-tree-open-delivery.md` / `.json`

### Wukong-only Visible Path Sample
- `agoal`
- `agoal contract`
- `agoal contract detail`
- `agoal contract fields`
- `agoal contract list`
- `agoal contract update`
- `agoal report`
- `agoal report list-statistics`
- `agoal report submit-detail`
- `agoal scorecard`
- `agoal scorecard detail`
- `agoal scorecard entity-detail`
- `agoal scorecard update`
- `agoal strategy`
- `agoal strategy detail`
- `agoal strategy list`
- `agoal strategy update`
- `agoal user`
- `agoal user objectives`
- `agoal user rules`
- `aiapp`
- `aiapp create`
- `aiapp modify`
- `aiapp query`
- `aidesign`
- `aidesign edit`
- `aidesign generate`
- `aidesign generate-with-image`
- `aidesign generate-with-template`
- `aidesign isolate`
- `aidesign upscale`
- `blackboard`
- `blackboard create`
- `blackboard list`
- `conference`
- `conference meeting`
- `conference meeting reserve`
- `conference member`
- `conference member invite`
- `contact label`
- `contact label get`
- `contact label list`
- `contact label list-members`
- `credit`
- `credit annual`
- `credit bidding`
- `credit branch`
- `credit cert-info`
- `credit change`
- `credit equity`
- `credit equity invest`
- `credit equity shareholder`
- `credit info`
- `credit ip`
- `credit ip copyright`
- `credit ip icp`
- `credit ip patent`
- `credit ip trademark`
- `credit kp`
- `credit license`
- `credit member`
- `credit risk`
- `credit risk assist`
- `credit risk consum`
- `credit risk court`
- `credit risk dishonest`
- `credit risk execute`
- `credit risk finalcase`
- `credit risk litigation`
- `credit risk owetax`
- `credit risk penalty`
- `credit risk pledge`
- `credit risk taxviolation`
- `credit risk verdict`
- `credit search`
- `docparse`
- `docparse convert`
- `finance`
- `finance account`
- `finance account list`
- `finance bank`
- `finance bank create`
- `finance bank list`
- `finance bank query`
- `finance category`
- `finance category search`
- `finance company`
- `finance company save`
- `finance company search`
- `finance company update`
- `finance customer`
- `finance customer get`
- `finance customer list`
- `finance customer save`
- `finance digital-invoice`
- `finance digital-invoice account`
- `finance digital-invoice batch-draw`
- `finance digital-invoice batch-draw-query`
- `finance digital-invoice batch-draw-query-saas`
- `finance digital-invoice batch-draw-saas`
- `finance digital-invoice do-login`
- `finance digital-invoice do-login-status`
- `finance digital-invoice face-qr`
- `finance digital-invoice face-status`
- `finance digital-invoice file`
- `finance digital-invoice get-table`
- `finance digital-invoice goods-code`
- `finance digital-invoice import-goods`
- `finance digital-invoice issue`
- `finance digital-invoice login-page`
- `finance digital-invoice search-goods`
- `finance digital-invoice send-email`
- `finance digital-invoice send-email-saas`
- `finance digital-invoice skill-version`
- `finance digital-invoice sms-code`
- `finance digital-invoice title`
- `finance gather`
- `finance gather execute`
- `finance gather execute-general`
- `finance gather query-rule`
- ... 80 more

### Open-only Visible Path Sample
- `chat media`
- `chat media upload`
- `doctor`
- `pat browser-policy`
- `skill get`
- `upgrade`

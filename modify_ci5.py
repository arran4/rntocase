import sys
import re

with open('.github/workflows/ci.yaml', 'r') as f:
    content = f.read()

# 1. Update workflow_dispatch inputs
inputs_old = r'''  workflow_dispatch:
    inputs:
      mode:
        description: "Operation mode"
        type: choice
        options:
          - lint-fix
          - publish-tag
          - release-major
          - release-minor
          - release-patch
          - release-test
          - release-rc
          - release-alpha
          - monthly-maintenance'''

inputs_new = r'''  workflow_dispatch:
    inputs:
      mode:
        description: "Operation mode"
        type: choice
        required: true
        default: build
        options: [build, lint-fix, monthly-maintenance, release-major, release-minor, release-patch, release-test, release-rc, release-alpha, publish-tag]'''

if inputs_old in content:
    content = content.replace(inputs_old, inputs_new)
else:
    # Handle potentially slightly different formatting if the above doesn't match exactly
    content = re.sub(
        r'  workflow_dispatch:\n    inputs:\n      mode:\n        description: "Operation mode"\n        type: choice\n        options:.*?- monthly-maintenance',
        r'  workflow_dispatch:\n    inputs:\n      mode:\n        description: "Operation mode"\n        type: choice\n        required: true\n        default: build\n        options: [build, lint-fix, monthly-maintenance, release-major, release-minor, release-patch, release-test, release-rc, release-alpha, publish-tag]',
        content, flags=re.DOTALL
    )

# 2. Fix the env block in route to remove GITHUB_REF and set up script

route_env_script_old = re.compile(r'        env:\n          INPUT_MODE: \${{ github.event.inputs.mode }}\n          GITHUB_EVENT_NAME: \${{ github.event_name }}\n          GITHUB_REF_TYPE: \${{ github.ref_type }}\n          GITHUB_REF: \${{ github.ref }}\n        run: \|\n          set -euo pipefail\n\n          run_code_checks=true.*?echo "mode=\$mode" >> "\$GITHUB_OUTPUT"', re.DOTALL)

route_env_script_new = r'''        env:
          INPUT_MODE: ${{ github.event.inputs.mode }}
          EVENT_NAME: ${{ github.event_name }}
          REF_TYPE: ${{ github.ref_type }}
        run: |
          set -euo pipefail

          run_code_checks=false
          run_build=false
          run_release=false
          run_autofix=false
          run_publisher=false
          run_maintenance=false
          mode="${INPUT_MODE:-}"

          case "$EVENT_NAME" in
            pull_request)
              run_code_checks=true
              run_build=true
              ;;
            push)
              run_code_checks=true
              run_build=true
              if [[ "$REF_TYPE" == "tag" && "$GITHUB_REF" == refs/tags/v* ]]; then
                run_publisher=true
              fi
              ;;
            schedule)
              run_code_checks=true
              run_build=true
              run_maintenance=true
              ;;
            workflow_dispatch)
              case "$mode" in
                build)
                  run_code_checks=true
                  run_build=true
                  ;;
                lint-fix)
                  run_autofix=true
                  ;;
                monthly-maintenance)
                  run_code_checks=true
                  run_build=true
                  run_maintenance=true
                  ;;
                release-major|release-minor|release-patch|release-test|release-rc|release-alpha)
                  run_code_checks=true
                  run_build=true
                  run_release=true
                  ;;
                publish-tag)
                  if [[ "$REF_TYPE" != "tag" || "$GITHUB_REF" != refs/tags/v* ]]; then
                    echo "publish-tag requires a v* tag ref; got $GITHUB_REF" >&2
                    sh -c "exit 1"
                  fi
                  run_code_checks=true
                  run_build=true
                  run_publisher=true
                  ;;
                *)
                  echo "Unsupported manual mode: $mode" >&2
                  sh -c "exit 1"
                  ;;
              esac
              ;;
            *)
              echo "Unsupported event: $EVENT_NAME" >&2
              sh -c "exit 1"
              ;;
          esac

          echo "run_code_checks=$run_code_checks" >> "$GITHUB_OUTPUT"
          echo "run_build=$run_build" >> "$GITHUB_OUTPUT"
          echo "run_release=$run_release" >> "$GITHUB_OUTPUT"
          echo "run_autofix=$run_autofix" >> "$GITHUB_OUTPUT"
          echo "run_publisher=$run_publisher" >> "$GITHUB_OUTPUT"
          echo "run_maintenance=$run_maintenance" >> "$GITHUB_OUTPUT"
          echo "mode=$mode" >> "$GITHUB_OUTPUT"'''

content = route_env_script_old.sub(route_env_script_new, content)

with open('.github/workflows/ci.yaml', 'w') as f:
    f.write(content)

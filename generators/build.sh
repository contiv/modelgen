#!/bin/bash

set -euo pipefail

cat <<EOF
package generators

var templates = map[string]string{
EOF

for file in templates/*.tmpl
do
cat <<EOF
"$(basename -s .tmpl $file)": \`
EOF

sed 's/`/` + "`" + `/g' $file

cat <<EOF
  \`,
EOF
done

cat <<EOF
}
EOF

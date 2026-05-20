#!/usr/bin/env bash
# nexusmind — traz novas versões do upstream Wei-Shaw/sub2api.
set -euo pipefail
cd "$(dirname "$0")/.."

if ! git remote | grep -qx upstream; then
  git remote add upstream https://github.com/Wei-Shaw/sub2api.git
fi

echo "==> git fetch upstream"
git fetch upstream

NEW=$(git rev-list --count main..upstream/main)
echo "==> $NEW commit(s) novos no upstream/main"
git --no-pager log --oneline main..upstream/main | head -40

if [ "$NEW" -eq 0 ]; then
  echo "Nada para atualizar."
  exit 0
fi

echo "==> git merge upstream/main"
echo "    Conflitos esperados só nos touchpoints de nexusmind/CUSTOMIZATIONS.md."
echo "    Resolva mantendo AMBOS os lados (edições aditivas)."
git merge upstream/main || {
  echo
  echo "Conflitos. Arquivos marcados com 'nexusmind' são os nossos touchpoints."
  echo "Após resolver: git merge --continue && bash nexusmind/verify.sh"
  exit 1
}

echo "==> Merge OK. Rode: bash nexusmind/verify.sh"

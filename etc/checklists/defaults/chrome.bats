load test_helper

setup() {
  defaults read com.google.chrome LastRunAppBundlePath || skip
}

@test "Chrome safe browsing enabled" {
  echo_url_for_description "$BATS_TEST_DESCRIPTION"
  local safebrowsing=$(cat "$HOME/Library/Application Support/Google/Chrome/Default/Preferences" | python -c "import json,sys;obj=json.load(sys.stdin);print obj['safebrowsing']['enabled'] if 'safebrowsing' in obj else 'True'")
  [ $? -eq 0 ] && [ "$safebrowsing" = "True" ]
}

load test_helper

@test "Safari blocks third-party cookies" {
  echo_url_for_description "$BATS_TEST_DESCRIPTION"
  run defaults read -app Safari WebKitStorageBlockingPolicy
  [ $status -ne 0 ] || [ "$output" -eq 1 ] || [ "$output" -eq 2 ]
}

@test "Safari blocks pop-up windows" {
  echo_url_for_description "$BATS_TEST_DESCRIPTION"
  run defaults read -app Safari WebKitJavaScriptCanOpenWindowsAutomatically
  [ $status -ne 0 ] || [ "$output" -eq 0 ]
}

@test "Safari warns when vising a fraudulent website" {
  echo_url_for_description "$BATS_TEST_DESCRIPTION"
  run defaults read -app Safari WarnAboutFraudulentWebsites
  [ $status -ne 0 ] || [ "$output" -eq 1 ]
}

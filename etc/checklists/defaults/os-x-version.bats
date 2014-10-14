load test_helper

@test "OS X version is Mavericks or newer" {
  echo_url_for_description "$BATS_TEST_DESCRIPTION"
  (( $(echo `sw_vers -productVersion` | cut -d. -f2) >= 9 ))
}

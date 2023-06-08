#Add debug for FSBL(optional)
XSCTH_BUILD_DEBUG = "1"
  
#Enable appropriate FSBL debug or compiler flags
YAML_COMPILER_FLAGS_append = " -DFSBL_DEBUG_INFO -DRSA_SUPPORT"

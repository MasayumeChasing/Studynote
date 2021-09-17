cmake_minimum_required(VERSION 2.6)
project(holo_mps)

set(JVMP_PROJECT_NAME holo_mps)

if(NOT CMAKE_BUILD_TYPE)
	set(CMAKE_BUILD_TYPE "Release")
endif(NOT CMAKE_BUILD_TYPE)

#set(CMAKE_CXX_STANDARD 11)
set(CMAKE_SKIP_RPATH  TRUE)
set(CMAKE_VERBOSE_MAKEFILE on)
add_compile_options(-std=c++11 -fstack-protector-all -fPIE -ftrapv -g)
SET(CMAKE_EXE_LINKER_FLAGS "${CMAKE_EXE_LINKER_FLAGS} -Wl,-z,relro,-z,now,-z,noexecstack -pie")
set(EXECUTABLE_OUTPUT_PATH ${CMAKE_BINARY_DIR}/bin)
set(LIBRARY_OUTPUT_PATH ${CMAKE_BINARY_DIR}/bin)
add_definitions(-D_FORTIFY_SOURCE=2)

set(CMAKE_CXX_FLAGS_DEBUG "$ENV{CXXFLAGS} -g")
set(CMAKE_CXX_FLAGS_RELEASE "$ENV{CXXFLAGS} -Os")

set(TARGET_DIR ${PROJECT_SOURCE_DIR}/build/target)
add_subdirectory(src/common ${TARGET_DIR}/src/common)
add_subdirectory(src ${TARGET_DIR}/src)

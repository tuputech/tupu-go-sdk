# Package base-recognition
> TUPU identification interface general module

## Overview

Package base-recognition implements access to TUPU recognition general interface

The user describes the file type to be identified by creating `DataInfo` object, creates a `Handler` object, uses its three methods(`Recognition()`, `RecognitionWithURL()`, `RecognitionWithPath()`) to access the recognition interface, and uses `ParseResult()` method to parse the results

## Example

Package `speech-sync` is one of the example
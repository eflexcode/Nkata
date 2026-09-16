// Top-level build file where you can add configuration options common to all sub-projects/modules.
plugins {
    alias(libs.plugins.android.application) apply false
    alias(libs.plugins.kotlin.compose) apply false
//    alias(libs.plugins.kotlin.kapt) apply false
    id("com.google.devtools.ksp") version "2.0.20-1.0.25" apply false
//    id("kotlin-kapt")
//    alias(libs.plugins.google.ksp) //apply false
//    alias(libs.plugins.kotlin.kapt) apply false
//    id("com.android.lagacy-kapt") apply false
}

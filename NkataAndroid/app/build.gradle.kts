plugins {
    alias(libs.plugins.android.application)
    alias(libs.plugins.kotlin.compose)
    id("com.google.devtools.ksp")
}

android {
    namespace = "com.ifeanyi.nkataandroid"
    compileSdk {
        version = release(36) {
            minorApiLevel = 1
        }
    }

    defaultConfig {
        applicationId = "com.ifeanyi.nkataandroid"
        minSdk = 24
        targetSdk = 36
        versionCode = 1
        versionName = "1.0"

        testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
    }

    buildTypes {
        release {
            optimization {
                enable = false
            }
        }
    }
    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_11
        targetCompatibility = JavaVersion.VERSION_11
    }
    buildFeatures {
        compose = true
    }
}

dependencies {
    implementation(platform(libs.androidx.compose.bom))
    implementation(libs.androidx.activity.compose)
    implementation(libs.androidx.compose.material3)
    implementation(libs.androidx.compose.ui)
    implementation(libs.androidx.compose.ui.graphics)
    implementation(libs.androidx.compose.ui.tooling.preview)
    implementation(libs.androidx.core.ktx)
    implementation(libs.androidx.lifecycle.runtime.ktx)
    testImplementation(libs.junit)
    androidTestImplementation(platform(libs.androidx.compose.bom))
    androidTestImplementation(libs.androidx.compose.ui.test.junit4)
    androidTestImplementation(libs.androidx.espresso.core)
    androidTestImplementation(libs.androidx.junit)
    debugImplementation(libs.androidx.compose.ui.test.manifest)
    debugImplementation(libs.androidx.compose.ui.tooling)
//    implementation(libs.androidx.compose.)
    implementation(libs.androidx.navigation.compose)
    implementation(libs.kotlinx.serialization.json)

    // Optional: Material 3 Window Size Class (For adaptive/responsive layouts)
    implementation(libs.androidx.compose.material3.window.size.class1)

    // Retrofit Gson Converter (for JSON parsing)
    implementation(libs.converter.gson)
    implementation(libs.retrofit)

    //mvvm
    implementation(libs.accompanist.systemuicontroller)
    implementation(libs.androidx.lifecycle.viewmodel.compose)
    implementation(libs.multidex)

    //img loading
    implementation(libs.coil.compose)
    implementation (libs.glide)
    ksp ("com.github.bumptech.glide:compiler:4.13.2")
//    kapt("com.github.bumptech.glide:compiler:4.13.2")
    //database
    implementation (libs.androidx.room.runtime)
    ksp (libs.androidx.room.compiler)

    // Kotlin Extensions and Coroutines support for Room
    implementation (libs.androidx.room.ktx)

    //color rendering
    implementation(libs.androidx.palette)
}
Pod::Spec.new do |s|
    s.name = 'Picovoice-iOS'
    s.module_name = 'Picovoice'
    s.version = '1.1.2'
    s.license = {:type => 'Apache 2.0', :file => 'LICENSE'}
    s.summary = 'iOS SDK for the Picovoice Offline Voice Recognition Platform'
    s.description = 
    <<-DESC
    Picovoice is an end-to-end platform for building voice products on your terms. It enables creating voice experiences similar to Alexa and Google. 
    But it entirely runs 100% on-device.

    Picovoice is:
        * Private: Everything is processed offline. Intrinsically HIPAA and GDPR compliant.
        * Reliable: Runs without needing constant connectivity.
        * Zero Latency: Edge-first architecture eliminates unpredictable network delay.
        * Accurate: Resilient to noise and reverberation. It outperforms cloud-based alternatives by wide margins *.
        * Cross-Platform: Design once, deploy anywhere. Build using familiar languages and frameworks.
    DESC
    s.homepage = 'https://github.com/Picovoice/picovoice/tree/master/sdk/ios'
    s.author = { 'Picovoice' => 'hello@picovoice.ai' }
    s.source = { :git => "https://github.com/Picovoice/picovoice.git", :tag => "Picovoice-iOS-v1.1.2"} 
    s.ios.deployment_target = '9.0'
    s.swift_version = '5.0'
    s.ios.framework = 'AVFoundation'
    s.source_files = 'sdk/ios/*.{swift}'
    s.dependency 'Porcupine-iOS', '~> 1.9.3'
    s.dependency 'Rhino-iOS', '~> 1.6.3'
  end
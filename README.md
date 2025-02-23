# Brain Flicker

Brain Flicker is an experimental visual application that allows users to create controlled flickering effects between two images with adjustable rates.

## ⚠️ HEALTH WARNING

**IMPORTANT: This application contains flashing images that may trigger seizures in people with photosensitive epilepsy.**

- DO NOT use this application if you have epilepsy, a history of seizures, or other photosensitive conditions
- DO NOT use this application if you are tired, sleep-deprived, or under the influence of alcohol
- Stop using immediately if you experience any discomfort, dizziness, nausea, disorientation, or visual disturbances
- Keep the room well lit when using this application
- Take regular breaks and limit usage time
- Maintain a safe distance from the screen

**Medical Disclaimer:** This application is for experimental purposes only. It is not a medical device and is not intended to diagnose, treat, cure, or prevent any disease or condition. Use at your own risk. The developers are not responsible for any adverse effects from using this application.

## Research Background

This application is inspired by research on cerebrospinal fluid flow during wakefulness in humans (April 2023). The implementation follows the experimental protocol where participants viewed a flickering checkerboard pattern with specific timing:

- Duration: 256 seconds (8 cycles)
- Cycle: 16 seconds on, 16 seconds off
- Pattern: Alternating checkerboard phases

Special thanks to Paul Keeble for the original research implementation and checkerboard pattern design.
Visit [Paul Keeble page](https://www.paulkeeble.co.uk/posts/cff/) for the original research context and implementation.

*Note: This is an experimental implementation based on published research. Results may vary.*

## Features

- Adjustable flicker rate between 1-99 flashes per second
- Simple interface with start/stop controls
- Rate adjustment during operation
- Image scaling that maintains aspect ratio
- Real-time visual feedback

## Usage

1. Launch the application
2. Set your desired flicker rate using the number input (1-99)
3. Click "Start" to begin the flicker effect
4. Use "Stop" to halt the effect
5. Click "Set" to change the rate while running
6. Access additional information via the "About" button

## Technical Requirements

- Operating System: Windows, or Linux
- Display: Monitor capable of displaying at least 60Hz refresh rate

## Installation
Download from release page appropriate version for your system version and just run.

## Planed
- [ ] Usage history
- [ ] Custom plan
- [ ] Option for sending reports
- [ ] Link in about goes to page with link to studies about this and more info


## Credits

https://www.paulkeeble.co.uk/posts/cff/ for inspiration
https://gioui.org/ great and efficent GUI library for develompent GOLANG apps

## License
MIT License, full text in LICENSE.md
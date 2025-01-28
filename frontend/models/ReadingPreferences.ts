export interface ReadingPreferences {
    atomicReading: boolean,
    font: string,
    theme: string,
    tts: {
        autoplay: boolean,
        voice: string,
        speed: number,
    }
}
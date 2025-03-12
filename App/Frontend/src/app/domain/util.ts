export function vhToPixels(vh: number): number {
    return (vh * window.innerHeight) / 100;
}

export function vwToPixels(vw: number): number {
    return (vw * window.innerWidth) / 100;
}
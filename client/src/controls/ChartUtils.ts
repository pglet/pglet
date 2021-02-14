export function parseNumber(n: any): number {
    try {
        const v = parseFloat(n.toString());
        return isNaN(v) ? 0 : v;
    } catch {
        return 0;
    }
}
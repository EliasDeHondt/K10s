
export function showTooltip(event: MouseEvent, data: Record<string, any> | undefined, tooltipId = 'tooltip') {
    const tooltip = document.getElementById(tooltipId);
    if (tooltip) {
        tooltip.textContent = data && typeof data === 'object'
            ? Object.entries(data)
                .map(([key, value]) => `${key}: ${value}`)
                .join('\n')
            : 'No data available';

        tooltip.style.display = 'block';
        tooltip.style.left = `${event.pageX + 10}px`;
        tooltip.style.top = `${event.pageY + 10}px`;
    }
}

export function hideTooltip(tooltipId = 'tooltip') {
    const tooltip = document.getElementById(tooltipId);
    if (tooltip) {
        tooltip.style.display = 'none';
    }
}
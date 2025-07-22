import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';
import type { TransitionConfig } from 'svelte/transition';
import { cubicOut, elasticOut } from 'svelte/easing';

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs));
}

export function ParseParams(url: string) {
    if (!url) return {};

    try {
        const params = new URLSearchParams(new URL(url).search);
        return Object.fromEntries(params);
    } catch (error) {
        console.error('Failed to ParseParams', url, error);
        return {};
    }
}

type FlyAndScaleParams = {
    y?: number;
    x?: number;
    start?: number;
    duration?: number;
    delay?: number;
};

export const flyAndScale = (
    node: Element,
    params: FlyAndScaleParams = { y: -8, x: 0, start: 0.95, duration: 150, delay: 0 }
): TransitionConfig => {
    const style = getComputedStyle(node);
    const transform = style.transform === 'none' ? '' : style.transform;

    const scaleConversion = (valueA: number, scaleA: [number, number], scaleB: [number, number]) => {
        const [minA, maxA] = scaleA;
        const [minB, maxB] = scaleB;

        const percentage = (valueA - minA) / (maxA - minA);
        const valueB = percentage * (maxB - minB) + minB;

        return valueB;
    };

    const styleToString = (style: Record<string, number | string | undefined>): string => {
        return Object.keys(style).reduce((str, key) => {
            if (style[key] === undefined) return str;
            return str + key + ':' + style[key] + ';';
        }, '');
    };

    return {
        duration: params.duration ?? 200,
        delay: params.delay ?? 0,
        css: (t) => {
            const y = scaleConversion(t, [0, 1], [params.y ?? 5, 0]);
            const x = scaleConversion(t, [0, 1], [params.x ?? 0, 0]);
            const scale = scaleConversion(t, [0, 1], [params.start ?? 0.95, 1]);

            return styleToString({
                transform: transform + 'translate3d(' + x + 'px, ' + y + 'px, 0) scale(' + scale + ')',
                opacity: t
            });
        },
        easing: cubicOut
    };
};

export function Pop(node: Element, { duration = 1000, delay = 0, power = 1, easing = elasticOut }) {
    return {
        duration,
        delay,
        easing,
        css: (t: number) => {
            const scale = 1 + Math.sin(t * Math.PI) * power;
            return `transform: scale(${scale});`;
        }
    };
}

export function Rotate(node: Element, { duration = 1000, delay = 0, easing = elasticOut }) {
    return {
        duration,
        delay,
        easing,
        css: (t: number) => {
            const angle = t * 360;
            const scale = 1 + Math.sin(t * Math.PI) * 0.05;
            return `transform: rotate(${angle}deg) scale(${scale});`;
        }
    };
}

export function Reveal(
    node: Element,
    { duration = 4500, easing = cubicOut, baseSpeed = 600 } = {}
) {
    const originalDisplay = window.getComputedStyle(node).display;
    const wrapper = document.createElement('div');
    wrapper.style.position = 'relative';
    wrapper.style.display = originalDisplay;
    node.parentNode?.insertBefore(wrapper, node);
    wrapper.appendChild(node);

    const animationLayer = document.createElement('div');
    animationLayer.style.position = 'absolute';
    animationLayer.style.top = '0';
    animationLayer.style.left = '0';
    animationLayer.style.width = '100%';
    animationLayer.style.height = '100%';
    wrapper.insertBefore(animationLayer, node);

    let spans: HTMLElement[] = [];
    let isAnimationComplete = false;

    function processNode(element: Node, parentElement: Element) {
        if (element.nodeType === Node.TEXT_NODE) {
            const words = element.textContent!.split(' ');
            words.forEach((word, wordIndex) => {
                const wordSpan = document.createElement('span');
                wordSpan.style.display = 'inline-block';
                wordSpan.style.whiteSpace = 'nowrap';

                word.split('').forEach((char) => {
                    const span = document.createElement('span');
                    span.textContent = char;
                    span.style.opacity = '0';
                    span.style.transform = 'translateY(20px)';
                    span.style.display = 'inline-block';
                    span.setAttribute('aria-hidden', 'true');

                    wordSpan.appendChild(span);
                    spans.push(span);
                });

                parentElement.appendChild(wordSpan);

                if (wordIndex < words.length - 1) {
                    const spaceSpan = document.createElement('span');
                    spaceSpan.innerHTML = '&nbsp;';
                    parentElement.appendChild(spaceSpan);
                }
            });
        } else if (element.nodeType === Node.ELEMENT_NODE) {
            const originalElement = element as HTMLElement;
            const newElement = originalElement.cloneNode(false) as HTMLElement;

            Array.from(originalElement.attributes).forEach((attr) => {
                newElement.setAttribute(attr.name, attr.value);
            });

            parentElement.appendChild(newElement);

            Array.from(originalElement.childNodes).forEach((child) => {
                processNode(child, newElement);
            });
        }
    }

    const clone = node.cloneNode(true);
    processNode(clone, animationLayer);

    // Hide original content during animation
    (node as HTMLElement).style.opacity = '0';

    const totalChars = spans.length;
    const stagger = totalChars > 1 ? (duration - baseSpeed) / (totalChars - 1) : 0;

    return {
        duration,
        tick: (t: number) => {
            const elapsed = t * duration;

            spans.forEach((span, i) => {
                const charDelay = i * stagger;
                const charAnimDuration = duration - charDelay;
                let progress = 0;

                if (elapsed >= charDelay) {
                    progress = Math.min(1, (elapsed - charDelay) / charAnimDuration);
                    progress = easing(progress);
                }

                span.style.opacity = progress.toString();
                span.style.transform = `translateY(${(1 - progress) * 20}px)`;
            });

            if (t === 1 && !isAnimationComplete) {
                isAnimationComplete = true;
                requestAnimationFrame(() => {
                    // Show original content
                    (node as HTMLElement).style.opacity = '1';
                    // Remove animation layer
                    animationLayer.remove();
                    // Unwrap the node
                    wrapper.parentNode?.insertBefore(node, wrapper);
                    wrapper.remove();
                });
            }
        }
    };
}

export const focus = (node: HTMLElement) => {
    node.focus();
};
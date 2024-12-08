/**********************************/
/* @since 01/01/2025             */
/* @author K10s Open Source Team */
/**********************************/

// Background animations for login page.
document.querySelectorAll('g').forEach(function(cube, index) {
    var currentTranslateY = parseFloat(cube.getAttribute('transform').split('(')[1].split(')')[0].split(',')[1]) || 0;
    var currentTranslateX = parseFloat(cube.getAttribute('transform').split('(')[1].split(')')[0].split(',')[0]) || 0;
    var scale = parseFloat(cube.getAttribute('transform').split('scale(')[1].split(')')[0]) || 1;

    anime({
        targets: cube,
        translateY: [currentTranslateY, currentTranslateY - 150],
        translateX: [currentTranslateX, currentTranslateX], // No animation
        scale: [scale, scale], // No animation
        duration: 1500, // 1.5 seconds
        direction: 'alternate',
        loop: true,
        delay: index * 100,
        endDelay: (el, i, l) => (l - i) * 100
    });
});

// Logo animation for login page.
document.getElementById('logo').onload = function() {
    var logo = document.getElementById('logo');
    setTimeout(function() {
        logo.src = '/images/icon-animation.png';
    }, 3000);
};

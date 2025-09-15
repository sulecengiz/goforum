/*!
* Start Bootstrap - Clean Blog v5.1.0 (https://startbootstrap.com/theme/clean-blog)
* Copyright 2013-2021 Start Bootstrap
* Licensed under MIT (https://github.com/StartBootstrap/startbootstrap-clean-blog/blob/master/LICENSE)
*/
(function ($) {
    "use strict"; // Start of use strict

    // Floating label headings for the contact form
    $("body").on("input propertychange", ".floating-label-form-group", function (e) {
        $(this).toggleClass("floating-label-form-group-with-value", !!$(e.target).val());
    }).on("focus", ".floating-label-form-group", function () {
        $(this).addClass("floating-label-form-group-with-focus");
    }).on("blur", ".floating-label-form-group", function () {
        $(this).removeClass("floating-label-form-group-with-focus");
    });

    // Show the navbar when the page is scrolled up
    var MQL = 992;

    //primary navigation slide-in effect
    if ($(window).width() > MQL) {
        var headerHeight = $('#mainNav').height();
        $(window).on('scroll', {
                previousTop: 0
            },
            function () {
                var currentTop = $(window).scrollTop();
                //check if user is scrolling up
                if (currentTop < this.previousTop) {
                    //if scrolling up...
                    if (currentTop > 0 && $('#mainNav').hasClass('is-fixed')) {
                        $('#mainNav').addClass('is-visible');
                    } else {
                        $('#mainNav').removeClass('is-visible is-fixed');
                    }
                } else if (currentTop > this.previousTop) {
                    //if scrolling down...
                    $('#mainNav').removeClass('is-visible');
                    if (currentTop > headerHeight && !$('#mainNav').hasClass('is-fixed')) $('#mainNav').addClass('is-fixed');
                }
                this.previousTop = currentTop;
            });
    }

    // Like (beğeni) toggling
    $(document).on('click', '.like-btn', function (e) {
        e.preventDefault();
        var $btn = $(this);
        if ($btn.data('busy')) return; // Spam engelle
        var commentId = $btn.data('comment-id');
        if (!commentId) return;
        $btn.data('busy', true);
        fetch('/like-comment/' + commentId, {
            method: 'POST',
            headers: { 'X-Requested-With': 'XMLHttpRequest' }
        }).then(function (res) {
            if (res.status === 401) { window.location.href = '/login'; return Promise.reject(); }
            if (!res.ok) return Promise.reject();
            return res.json();
        }).then(function (data) {
            if (!data || !data.success) return;
            var likeCountEl = document.getElementById('like-count-' + commentId);
            var iconEl = document.getElementById('like-icon-' + commentId);
            var labelEl = document.getElementById('like-label-' + commentId);
            if (likeCountEl) likeCountEl.textContent = data.likeCount;
            if (iconEl) {
                if (data.isLiked) { iconEl.classList.remove('far'); iconEl.classList.add('fas'); $btn.addClass('liked'); $btn.attr('aria-pressed','true'); }
                else { iconEl.classList.remove('fas'); iconEl.classList.add('far'); $btn.removeClass('liked'); $btn.attr('aria-pressed','false'); }
            }
            if (labelEl) {
                if (data.likeCount === 0) labelEl.textContent = 'Beğen';
                else if (data.likeCount === 1) labelEl.textContent = '1 beğeni';
                else labelEl.textContent = data.likeCount + ' beğeni';
            }
        }).catch(function () { /* sessiz geç */ }).finally(function () { $btn.data('busy', false); });
    });

})(jQuery); // End of use strict
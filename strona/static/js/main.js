function disable($form) {
  $form.find('input').attr('disabled', true);
  $form.find('textarea').attr('disabled', true);
}

function enable($form) {
  $form.find('input').attr('disabled', false);
  $form.find('textarea').attr('disabled', false);
}

function message(text, type) {
  var $m = $('#message');
  if ('error' === type) {
    $m.text(text);
    $m.addClass('text-danger');

    return;
  }

  if ('info' === type) {
    $m.text(text);
    $m.addClass('text-info');

    return;
  }

  $m.text('');
  $m.removeClass('text-danger');
  $m.removeClass('text-info');
}

function loadCaptcha(API_URL, $form) {
  $.ajax({
    type: "POST",
    url: API_URL + "/captchas",
    success: function(data) {
      $('#captchaImg').attr('src', 'data:image/png;base64,' + data.captcha);
      _CAPTCHA = data.id;

      enable($form);
    }
  });
}

var _CAPTCHA = null;

$(document).ready(function(e) {
  var $form = $('#contactForm');
  var API_URL = $form.attr('action');
  disable($form);
  if ('/contact/' === window.location.pathname) {
    loadCaptcha(API_URL, $form);
  }

  $('#formRegenerateCaptcha').on('click', function(e) {
    e.preventDefault();
    message();
    loadCaptcha(API_URL, $form);
  });

  $('#contactForm').on('submit', function(e) {
    e.preventDefault();
    disable($form);

    var req = {
      data: {
        title: $('#title').val(),
        content: $('#content').val()
      },
      captcha: {
        id: _CAPTCHA,
        secret: $('#captcha').val()
      }
    };

    $.ajax({
      type: "POST",
      url: API_URL + "/forms",
      dataType: 'json',
      contentType: 'application/json',
      data: JSON.stringify(req),
      success: function(data) {
        enable($form);

        if (data.message !== 'Yay') {
          message(data.message, 'error');
          return;
        }
        message(data.message, 'info');
      },
      error: function() {
        console.log('Buu...');
        enable($form);
        message('Buu...', 'error');
      }
    });
  });
});

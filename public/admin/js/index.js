$(document).ready(function () {

  // 如果是写文章页面 暂时隐藏
  if (HIDE_SIDE_BAR) {
    hideBar();
  }
  // content 重重新计算高度
  contentResize();

  // 日期选择插件
  $('#datetimepicker').datetimepicker({});

  // 模拟滚动条
  $('#sidebar').niceScroll({cursorcolor: 'rgba(0,0,0,0.3)'});

  // 点击菜单按钮侧边栏
  $('#sidebar-control').click(function () {
  $('body').toggleClass('sidebar-toggle');
  });
  $('.sub-menu').click(function () {
  var ul = $(this).children('ul');
  var height = ul.height();
  console.log(height);
  if (ul.css('display') == 'none') {
    $(this).addClass('toggle');
    ul.css({display: 'block', height: 0});
    ul.animate({'height': height}, 300);

  } else {
    $(this).removeClass('toggle');
    ul.animate({'height': 0}, 300, function () {
    ul.css({display: 'none', height: height});
    });
  }
  });

  // 文章右侧选项切换
  $('.post-options>ul>li').click(function () {
  if ($(this).hasClass('active')) {
    return;
  }
  $(this).siblings().removeClass('active');
  $(this).addClass('active');

  $('#post-option-cont').toggle();
  $('#post-attachment-cont').toggle();
  });

  $('.tags-content>span').click(function () {
  if ($(this).hasClass('active')) {
    $(this).removeClass('active');
  } else {
    $(this).addClass('active');
  }
  });

  // 标签X掉
  $('.post-tags>ul').on('click', 'li>span', function () {
  $(this).parent('.tag-li').remove();
  });

  $('#tag-input').on('keyup', function (e) {
  var str = $(this).val().trim();
  if (str.length == 0) {
    return;
  }
  if (e.keyCode == 13) {
    var temp = '<li class="tag-li">'+str+' <span>x</span></li>';
    $(this).parent().before(temp);
    $(this).val('');
  }
  });
  
  // 提交文章
  $('#publish').click(function () {
    var $btn = $(this).button('loading');
    var post = {};
    post.title = $('#post-title').val();
    post.content = $('#post-content').val();
    post.type = 'post';
    post.status = 'publish';
    $.post('/admin/post', post, function (result) {
      if (result.code == 200) {
        pop(result.msg);
        result.redirect && (window.location.href = result.redirect);
        
      } else {
        pop(result.msg);
        result.redirect && (window.location.href = result.redirect);
      }
      $btn.button('reset');
    });
  })

});

var contentResize = function () {
  var wh = $(window).height();
  var h = $('.content').height();
  var th = wh - 58 - 70 - 20;
  if (h < th) {
  $('.content').css({'min-height': th});
  }
}

var hideBar = function () {
  $('body').addClass('hide-sidebar');
}

var pop = function (text) {
  $('.overlay').clearQueue().text(text || '').
      stop().
      fadeIn().
      delay(3000).
      fadeOut();
};
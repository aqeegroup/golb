$(document).ready(function () {

  // 如果是写文章页面 暂时隐藏
  if (HIDE_SIDE_BAR) {
    hideBar();
  }
  // content 重重新计算高度
  contentResize();

  // 日期选择插件
  $('#datetimepicker').datetimepicker({
    useCurrent: false,
  });

  // 模拟滚动条
  $('#sidebar').niceScroll({ cursorcolor: 'rgba(0,0,0,0.3)' });

  // 点击菜单按钮侧边栏
  $('#sidebar-control').click(function () {
    $('body').toggleClass('sidebar-toggle');
  });
  $('.sub-menu').click(function () {
    var ul = $(this).children('ul');
    var height = ul.height();
    // console.log(height);
    if (ul.css('display') == 'none') {
      $(this).addClass('toggle');
      ul.css({ display: 'block', height: 0 });
      ul.animate({ 'height': height }, 300);

    } else {
      $(this).removeClass('toggle');
      ul.animate({ 'height': 0 }, 300, function () {
        ul.css({ display: 'none', height: height });
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
      var temp = '<li class="tag-li">' + str + ' <span>x</span></li>';
      $(this).parent().before(temp);
      $(this).val('');
    }
  });

  // 提交文章
  $('#publish').click(function () {
    var btn = $(this).button('loading');
    var post = {};
    post.title = $('#post-title').val();
    post.content = $('#post-content').val();
    post.type = 'post';
    post.status = 'publish';
    post.slug = $('#slug').val()

    var cateSelected = [];
    $('.cate input:checked').each(function() {
      cateSelected.push($(this).val());
    });

    post.cates = cateSelected.join(',');

    $.post('/admin/post', post, function (data) {
      pop(data.msg);      
      if (data.code == 200) {
        setTimeout(function () {
          data.redirect && (window.location.href = data.redirect);                              
        }, 1000);  
      }
      btn.button('reset');
    });
  });

  // 文章列表页全选功能
  $('#select-all').click(function () {
    var isChecked = $(this).prop("checked");
    $('#post-table').find('input').prop("checked", isChecked);
    isChecked && $('#post-table>tbody tr').addClass('checked');
    isChecked || $('#post-table>tbody tr').removeClass('checked');
  });
  // 单个文章点击tr选中
  $('#post-table>tbody>tr').click(function () {
    if ($(this).hasClass('checked')) {
      $(this).removeClass('checked');
      $(this).find('input').prop("checked", false);      
    } else {
      $(this).addClass('checked');
      $(this).find('input').prop("checked", true);
    }
  });
  // 选中项操作


  // 新增分类
  $('#submit-cate').click(function () {
    var modal = $('#add-cate');
    var form = {};
    form.name = modal.find('input[name=name]').val();
    form.slug = modal.find('input[name=slug]').val();
    form.parent_id = modal.find('select[name=parent]').val();
   
    // console.log(form);
    $.post('/admin/cate', form, function (data) {
      if (data.code == "200") {
        pop(data.msg);
        hideModal();
        data.redirect && (window.location.href = data.redirect);        
      } else {
        pop(data.msg);
      }
    });
  });

  // 删除分类
  $('#delete-cate').click(function () {

    Util.confirm('确认删除选中的分类吗？', function () {
      var checked = [];
      $('#cates').find('input:checked').each(function () {
        checked.push($(this).val());
      });

      // console.log(checked);
      $.post('/admin/cate/del', {ids: checked.join(',')}, function(data) {
        pop(data.msg);
        if (data.code == '200') {
          Util.close();
          setTimeout(function () {
            data.redirect && (window.location.href = data.redirect);                              
          }, 1000);            
        }
      });
    });
  });

  // 删除文章
  $('#delete-post').click(function () {
    Util.confirm('确认删除选中的文章吗？', function () {
      var checked = [];
      $('#post-list').find('input:checked').each(function () {
        checked.push($(this).val());
      });
      console.log(checked);

      $.post('/admin/post/del', {ids: checked.join(',')}, function(data) {
        pop(data.msg);        
        if (data.code == '200') {
          Util.close();
          setTimeout(function () {
            data.redirect && (window.location.href = data.redirect);                              
          }, 1000);
        }
      });
    });
  });
  

});

var contentResize = function () {
  var wh = $(window).height();
  var h = $('.content').height();
  var th = wh - 58 - 70 - 20;
  if (h < th) {
    $('.content').css({ 'min-height': th });
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

var showModal = function (id) {
  $('#'+id).show();
  $('#modal').modal('show');
}

var hideModal = function () {
  $(".modal-content").hide();
  $('#modal').modal('hide');
}
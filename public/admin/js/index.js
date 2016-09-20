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

  var clicking = false;
  $('.sub-menu').click(function () {
    if (clicking) return;
    clicking = true;
    var ul = $(this).children('ul');
    var height = ul.height();
    // console.log(height);
    if (ul.css('display') == 'none') {
      $(this).addClass('toggle');
      ul.css({ display: 'block', height: 0 });
      ul.animate({ 'height': height }, 300, function() {
        clicking = false;        
      });
    } else {
      $(this).removeClass('toggle');
      ul.animate({ 'height': 0 }, 300, function () {
        ul.css({ display: 'none', height: height });
        clicking = false;
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
    post.id = $('#id').val();
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
          btn.button('reset');
          data.redirect && (window.location.href = data.redirect);                              
        }, 1000);  
      }
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
    $('#cate-title').text('添加分类');
    var modal = $('#add-cate');
    var form = {};
    form.id = modal.find('input[name=id]').val();
    form.name = modal.find('input[name=name]').val();
    form.slug = modal.find('input[name=slug]').val();
    form.parent_id = modal.find('select[name=parent]').val();
   
    // console.log(form);
    $.post('/admin/cate', form, function (data) {
      if (data.code == "200") {
        pop(data.msg);
        hideModal();
        setTimeout(function () {
          data.redirect && (window.location.href = data.redirect);                              
        }, 1000);         
      } else {
        pop(data.msg);
      }
    });
  });

  // 修改分类
  $('#cates a').click(function () {
    var modal = $('#add-cate');    
    $('#cate-title').text('修改分类');
    var cate = $(this).data('cate');
    // console.log(cate);
    modal.find('input[name=id]').val(cate.ID);
    modal.find('input[name=name]').val(cate.Name);
    modal.find('input[name=slug]').val(cate.Slug);
    modal.find('select[name=parent]').val(cate.ParentID);
    showModal('add-cate');
    return false;
  });

  // 删除分类
  $('#delete-cate').click(function () {
    var checked = [];
    $('#cates').find('input:checked').each(function () {
      checked.push($(this).val());
    });
    if (checked.length == 0) {
      pop('请选择要删除的分类');
      return;
    }
    Util.confirm('确认删除选中的分类吗？', function () {
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
    var checked = [];
    $('#post-list').find('input:checked').each(function () {
      checked.push($(this).val());
    });
    if (checked.length == 0) {
      pop('请选择要删除的文章');
      return;
    }
    Util.confirm('确认删除选中的文章吗？', function () {
      // console.log(checked);
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
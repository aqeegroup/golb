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
      var temp = '<li class="tag-li">' + str + '<span><i class="iconfont icon-close"></i></span></li>';
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

    var tagSelected = [];
    $('.tags-content>.active').each(function(){
      tagSelected.push($(this).text().trim());
    });
    $('.post-tags .tag-li').each(function(){
      var tagName = $(this).text().trim();
      if (tagSelected.indexOf(tagName) < 0) {
        tagSelected.push(tagName);
      }
    });
    post.tags = tagSelected.join(',');
    // console.log(tagSelected);

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
    var cate = $(this).data('cate');
    // console.log(cate);
    modal.find('input[name=id]').val(cate.ID);
    modal.find('input[name=name]').val(cate.Name);
    modal.find('input[name=slug]').val(cate.Slug);
    modal.find('select[name=parent]').val(cate.ParentID);
    showModal('add-cate', '修改分类');
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

  // 选中标签
  $('.tag-manage li').click(function() {
    $(this).toggleClass('active');
  });
  // 小铅笔点击弹出编辑
  $('.tag-manage li i').click(function() {
    var modal = $('#tag-edit');
    var tag = $(this).data('tag');
    // console.log(cate);
    modal.find('input[name=id]').val(tag.ID);
    modal.find('input[name=name]').val(tag.Name);
    modal.find('input[name=slug]').val(tag.Slug);
    showModal('tag-edit', '修改标签');
    return false;
  });
  // 提交标签编辑结果 AND 新建标签
  $('#submit-tag').click(function() {
    var modal = $('#tag-edit');
    var form = {};
    form.id = modal.find('input[name=id]').val();
    form.name = modal.find('input[name=name]').val();
    form.slug = modal.find('input[name=slug]').val();
    // console.log(form);
    $.post('/admin/tag', form, function(data) {
      pop(data.msg);
      if (data.code == '200') {
        Util.close();
        setTimeout(function () {
          data.redirect && (window.location.href = data.redirect);                              
        }, 1000);            
      }
    });

  });
  // delete-tag
  $('#delete-tag').click(function() {
    var checked = [];
    $('.tag-manage ul .active').each(function () {
      checked.push($(this).data('id'));
    });
    if (checked.length == 0) {
      pop('请选择要删除的标签');
      return;
    }
    Util.confirm('确认删除选中的标签吗？', function () {
      console.log(checked);
      $.post('/admin/tag/del', {ids: checked.join(',')}, function(data) {
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


  $('#modal').on('hide.bs.modal', function () {
    $(this).find("input").val('');
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

var showModal = function (id, title) {
  var modal = $('#'+id)
  title && modal.find('.modal-title').text(title)
  modal.show();
  $('#modal').modal('show');
}

var hideModal = function () {
  $(".modal-content").hide();
  $('#modal').modal('hide');
}

var clear = function () {

}

(function(window, $){
  window.Util = {};
  var Util = window.Util;
  

  function confirm(callback){

  }

  Util.alert = function(text) {
    $('#alert .modal-body').text(text);
    $('#alert').modal('show');
  }

  Util.confirm = function(text, callback) {
    $('#confirm .modal-body').text(text);
    callback && $('#confirm-btn').click(callback);
    $('#confirm').modal('show');    
  }

  Util.close = function() {
    $('#confirm').modal('hide');
  }

}(window, $))
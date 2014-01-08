var xrguide = angular.module('xrguide', ['ngRoute', 'xrguideControllers']);
xrguide.
    config(function ($routeProvider, $locationProvider) {
        'use strict';
        $routeProvider.
            when('/wares', {
                templateUrl: '/tmpl/wares.html',
                controller: 'WaresListCtrl'
            }).
            when('/', {
                templateUrl: '/tmpl/index.html'
            });
        $locationProvider.html5Mode(true);
    });

xrguide.service('Ware', ['$rootScope', function ($rootScope) {
    'use strict';
    var wares = [];
    var service = {
        wares: function() {
            return wares;
        },

        update: function ($http, $scope) {
            $http({
                method: 'GET',
                url: '/wares',
                headers: {
                    'Accept': 'application/json'
                }
            }).success(function (data) {
                $scope.wares = data;
            });
            $rootScope.$broadcast('wares.update');
        }
    };

    return service;
}]);


xrguide.directive('wareRow', function () {
    'use strict';
    return {
        restrict: 'E',
        scope: {
            ware: '='
        },
        'templateUrl': '/tmpl/ware.html'
    };
});

var xrguideControllers = angular.module('xrguideControllers', []);

xrguideControllers.controller('WaresListCtrl', ['Ware', '$scope', '$http', function (Ware, $scope, $http) {
    'use strict';
    Ware.update($http, $scope);
    $scope.internalClass = function (ware) {
        if (ware.Name.Valid !== undefined && ware.Name.Valid) {
            return '';
        }
        return 'ware-internal warning'
    };
    $scope.$watch('wares.update', function () {
    });
}]);


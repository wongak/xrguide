var xrguide = angular.module('xrguide', ['ngRoute', 'xrguideControllers']);
xrguide.
    config(function ($routeProvider, $locationProvider) {
        'use strict';
        $routeProvider.
            when('/wares', {
                templateUrl: '/tmpl/wares.html',
                controller: 'WaresListCtrl'
            }).
            when('/ware/:wareId', {
                templateUrl: '/tmpl/ware.html',
                controller: 'WareDetailCtrl'
            }).
            when('/', {
                templateUrl: '/tmpl/index.html'
            });
        $locationProvider.html5Mode(true);
    });

xrguide.service('Ware', ['$rootScope', '$http', function ($rootScope, $http) {
    'use strict';
    var wares,
        service,
        wareDefaults;

    wares = [];
    wareDefaults = {
        wareName: function (ware) {
            if (ware === undefined) {
                ware = this;
            }
            if (ware.Name.Valid !== undefined && ware.Name.Valid) {
                return ware.Name.String;
            }
            if (ware.NameRaw.Valid !== undefined && ware.NameRaw.Valid) {
                return ware.NameRaw.String;
            }
            return '?';
        },

        wareSpecialist: function (ware) {
            if (ware === undefined) {
                ware = this;
            }
            if (ware.Specialist.Valid !== undefined && ware.Specialist.Valid) {
                return ware.Specialist.String;
            }
            return '-';
        },

        wareHasProduction: function (ware) {
            if (ware === undefined) {
                ware = this;
            }
            if (ware.Productions === undefined) {
                return false;
            }
            return true;
        }
    };
    service = {
        wares: function () {
            return wares;
        },

        updateWares: function ($scope) {
            $http({
                method: 'GET',
                url: '/wares',
                headers: {
                    'Accept': 'application/json'
                }
            }).success(function (data) {
                $scope.wares = angular.extend(data, wareDefaults);
                $rootScope.$broadcast('wares.update');
            });
        },

        getWare: function (wareId, $scope) {
            $http({
                method: 'GET',
                url: '/ware/' + wareId,
                headers: {
                    'Accept': 'application/json'
                }
            }).success(function (data) {
                $scope.ware = angular.extend(data, wareDefaults);
                $scope.$broadcast('ware.update');
            });
        }
    };

    return service;
}]);

xrguide.directive('wareProduction', function () {
    return {
        restrict: 'E',
        templateUrl: '/tmpl/wareProduction.html',

    };
})

var xrguideControllers = angular.module('xrguideControllers', []);

xrguideControllers.controller('WaresListCtrl', ['Ware', '$scope', function (Ware, $scope) {
    'use strict';
    Ware.updateWares($scope);

    $scope.query = '';
    $scope.filterInternal = true;
    $scope.showTranspEquipment = false;
    $scope.showTranspInventory = false;
    $scope.showTranspEnergy = true;
    $scope.showTranspContainer = true;
    $scope.showTranspBulk = true;
    $scope.showTranspLiquid = true;
    $scope.showTranspFuel = true;

    $scope.isInternal = function (ware) {
        if (ware.Name.Valid !== undefined && ware.Name.Valid) {
            return false;
        }
        return true;
    };
    $scope.internalClass = function (ware) {
        if (!$scope.isInternal(ware)) {
            return '';
        }
        return 'ware-internal warning';
    };
    $scope.filterInternalExp = function (ware) {
        if (!$scope.filterInternal) {
            return true;
        }
        return !$scope.isInternal(ware);
    };
    $scope.filterTransports = function (ware) {
        if (!$scope.showTranspEquipment && ware.Transport === 'equipment') {
            return false;
        }
        if (!$scope.showTranspInventory && ware.Transport === 'inventory') {
            return false;
        }
        if (!$scope.showTranspEnergy && ware.Transport === 'energy') {
            return false;
        }
        if (!$scope.showTranspContainer && ware.Transport === 'container') {
            return false;
        }
        if (!$scope.showTranspBulk && ware.Transport === 'bulk') {
            return false;
        }
        if (!$scope.showTranspLiquid && ware.Transport === 'liquid') {
            return false;
        }
        if (!$scope.showTranspFuel && ware.Transport === 'fuel') {
            return false;
        }
        return true;
    };

//    $scope.$watch('wares.update', function () {});
}]);

xrguideControllers.controller('WareDetailCtrl', ['Ware', '$scope', '$routeParams', function (Ware, $scope, $routeParams) {
    'use strict';
    Ware.getWare($routeParams.wareId, $scope);
}]);
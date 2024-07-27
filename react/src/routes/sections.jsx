import { lazy, Suspense } from 'react';
import { Outlet, Navigate, useRoutes } from 'react-router-dom';

import DashboardLayout from 'src/layouts/dashboard';
import IssuesPage from 'src/pages/issues';
import PrivacyPage from 'src/pages/privacy-policy';
import TermsPage from 'src/pages/terms';

export const IndexPage = lazy(() => import('src/pages/app'));
export const BlogPage = lazy(() => import('src/pages/blog'));
export const SearchPage = lazy(() => import('src/pages/search'));
export const SearchUserPage = lazy(() => import('src/pages/search-user'));
export const LoginPage = lazy(() => import('src/pages/login'));
export const ProductsPage = lazy(() => import('src/pages/products'));
export const Page404 = lazy(() => import('src/pages/page-not-found'));

// ----------------------------------------------------------------------

export default function Router() {
    const routes = useRoutes([
        {
            element: (
                <DashboardLayout>
                    <Suspense>
                        <Outlet />
                    </Suspense>
                </DashboardLayout>
            ),
            children: [
                { element: <IndexPage />, index: true },
                { path: 'search-user', element: <SearchUserPage /> },
                { path: 'search', element: <SearchPage /> },
                { path: 'issues', element: <IssuesPage /> },
                { path: 'products', element: <ProductsPage /> },
                { path: 'blog', element: <BlogPage /> },
            ],
        },
        {
            path: 'login',
            element: <LoginPage />,
        },
        {
            path: 'privacy-policy',
            element: <PrivacyPage />,
        },
        {
            path: 'terms-and-conditions',
            element: <TermsPage />,
        },
        {
            path: '404',
            element: <Page404 />,
        },
        {
            path: '*',
            element: <Navigate to="/404" replace />,
        },
    ]);

    return routes;
}

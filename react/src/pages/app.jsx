import { Helmet } from 'react-helmet-async';

import { SearchUserView } from 'src/sections/search-user/view';

// ----------------------------------------------------------------------

export default function AppPage() {
    return (
        <>
            <Helmet>
                <title> User Search | GooGL-Bye </title>
            </Helmet>

            <SearchUserView />
        </>
    );
}
